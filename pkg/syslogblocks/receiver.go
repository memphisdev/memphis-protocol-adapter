package syslogblocks

import (
	"github.com/g41797/sputnik"
)

const (
	ReceiverName           = "syslogreceiver"
	ReceiverResponsibility = "syslogreceiver"

	ProducerName           = "syslogproducer"
	ProducerResponsibility = "syslogproducer"

	StorerName           = "syslogstorer"
	StorerResponsibility = "syslogstorer"
)

func ReceiverDescriptor() sputnik.BlockDescriptor {
	return sputnik.BlockDescriptor{Name: ReceiverName, Responsibility: ReceiverResponsibility}
}

func init() {
	sputnik.RegisterBlockFactory(ReceiverName, receiverBlockFactory)
}

func receiverBlockFactory() *sputnik.Block {
	receiver := new(receiver)
	block := sputnik.NewBlock(
		sputnik.WithInit(receiver.init),
		sputnik.WithRun(receiver.run),
		sputnik.WithFinish(receiver.finish),
		sputnik.WithOnConnect(receiver.brokerConnected),
		sputnik.WithOnDisconnect(receiver.brokerDisconnected),
	)
	return block
}

type receiver struct {
	conf     SyslogConfiguration
	syslogd  *Server
	producer sputnik.BlockCommunicator
	backup   sputnik.BlockCommunicator

	// Used for synchronization
	// between finish and run
	stop chan struct{}
	done chan struct{}
}

// Init
func (rcv *receiver) init(fact sputnik.ConfFactory) error {
	if err := fact(ReceiverName, &rcv.conf); err != nil {
		return err
	}

	syslogd := NewServer(rcv.conf)

	if err := syslogd.Init(); err != nil {
		return err
	}

	syslogd.SetupHandling(nil)

	rcv.syslogd = syslogd
	rcv.stop = make(chan struct{}, 1)

	return nil
}

// Finish:
func (rcv *receiver) finish(init bool) {
	if init {
		rcv.stopSyslog()
		return
	}

	close(rcv.stop) // Cancel Run

	<-rcv.done // Wait finish of Run
	return
}

// OnServerConnect:
func (rcv *receiver) brokerConnected(_ sputnik.ServerConnection) {
	rcv.syslogd.SetupHandling(rcv.producer)
	return
}

// OnServerDisconnect:
func (rcv *receiver) brokerDisconnected() {
	// Redirect to syslog backup block for save logs in the storage
	rcv.syslogd.SetupHandling(rcv.backup)
	return
}

// Run:
func (rcv *receiver) run(bc sputnik.BlockCommunicator) {

	err := rcv.syslogd.Start()
	if err != nil {
		panic(err)
	}

	defer rcv.stopSyslog()

	rcv.done = make(chan struct{})
	defer close(rcv.done)

	producer, exists := bc.Communicator(ProducerResponsibility)
	if !exists {
		panic("Syslog producer block does not exists")
	}

	rcv.producer = producer

	// Storer (backup/restore) block - optional
	// If does not exists, all logs will be discarded
	rcv.backup, _ = bc.Communicator(StorerResponsibility)

	<-rcv.stop

	return
}

func (rcv *receiver) stopSyslog() {

	if rcv == nil {
		return
	}

	if rcv.syslogd == nil {
		return
	}

	rcv.syslogd.SetupHandling(nil)
	rcv.syslogd.Finish()

	return
}
