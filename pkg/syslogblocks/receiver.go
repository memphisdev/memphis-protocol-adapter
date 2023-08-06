package syslogblocks

import (
	"github.com/g41797/sputnik"
)

const (
	ReceiverName           = "syslogreceiver"
	ReceiverResponsibility = "syslogreceiver"

	ProducerName           = "syslogproducer"
	ProducerResponsibility = "syslogproducer"
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
	// For now - disable handling
	// Next stage - redirect to syslog backup block for save logs in the storage
	rcv.syslogd.SetupHandling(rcv.backup)
	return
}

// Run:
func (rcv *receiver) run(bc sputnik.BlockCommunicator) {

	rcv.done = make(chan struct{})
	defer close(rcv.done)

	producer, exists := bc.Communicator(ProducerResponsibility)
	if !exists {
		panic("Syslog producer block does not exists")
	}

	rcv.producer = producer

	select {
	case <-rcv.stop:
		rcv.syslogd.SetupHandling(nil)
		rcv.syslogd.Finish()
		return
	}

	return
}
