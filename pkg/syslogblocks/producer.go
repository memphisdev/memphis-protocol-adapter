package syslogblocks

import (
	"github.com/g41797/sputnik"
)

func producerBlockFactory() *sputnik.Block {
	prd := new(producer)
	if mpf == nil {
		return nil
	}

	mp := mpf()
	if mp == nil {
		return nil
	}

	prd.mp = mp

	block := sputnik.NewBlock(
		sputnik.WithInit(prd.init),
		sputnik.WithRun(prd.run),
		sputnik.WithFinish(prd.finish),
		sputnik.WithOnConnect(prd.brokerConnected),
		sputnik.WithOnDisconnect(prd.brokerDisconnected),
		sputnik.WithOnMsg(prd.logReceived),
	)
	return block
}

func ProducerDescriptor() sputnik.BlockDescriptor {
	return sputnik.BlockDescriptor{Name: ProducerName, Responsibility: ProducerResponsibility}
}

func init() {
	sputnik.RegisterBlockFactory(ProducerName, producerBlockFactory)
}

type producer struct {
	mp        sputnik.MessageProducer
	connected bool
	conf      MsgPrdConfig
	cfact     sputnik.ConfFactory
	backup    sputnik.BlockCommunicator
	stop      chan struct{}
	done      chan struct{}
	conn      chan struct{}
	dscn      chan struct{}
	mlog      chan sputnik.Msg
}

// Init
func (prd *producer) init(fact sputnik.ConfFactory) error {
	if err := fact(MsgProducerConfigName, &prd.conf); err != nil {
		return err
	}
	prd.cfact = fact
	prd.stop = make(chan struct{}, 1)
	prd.done = make(chan struct{}, 1)
	prd.conn = make(chan struct{}, 1)
	prd.dscn = make(chan struct{}, 1)
	prd.mlog = make(chan sputnik.Msg, 1)

	return nil
}

// Finish:
func (prd *producer) finish(init bool) {
	if init {
		return
	}

	close(prd.stop) // Cancel Run

	<-prd.done // Wait finish of Run
	return
}

// OnServerConnect:
func (prd *producer) brokerConnected(_ sputnik.ServerConnection) {
	prd.conn <- struct{}{}
	return
}

// OnServerDisconnect:
func (prd *producer) brokerDisconnected() {
	prd.dscn <- struct{}{}
	return
}

// OnMsg:
func (prd *producer) logReceived(msg sputnik.Msg) {
	prd.mlog <- msg
	return
}

// Run
func (prd *producer) run(bc sputnik.BlockCommunicator) {

	prd.backup, _ = bc.Communicator(StorerResponsibility)

	defer close(prd.done)

loop:
	for {
		select {
		case <-prd.stop:
			break loop
		case <-prd.conn:
			{
				if err := prd.mp.Connect(prd.cfact); err == nil {
					prd.connected = true
				} else {
					prd.connected = false
				}
			}
		case <-prd.dscn:
			{
				if prd.connected {
					prd.mp.Disconnect()
					prd.connected = false
				}
			}
		case logmsg := <-prd.mlog:
			prd.processLog(logmsg)
		}
	}

	prd.mp.Disconnect()
	return
}

func (prd *producer) processLog(logmsg sputnik.Msg) {
	sendToBackup := prd.backup != nil
	from, exists := logmsg["from"]
	if exists && (from == ProducerResponsibility) {
		sendToBackup = false
	}

	if !prd.connected && sendToBackup {
		logmsg["from"] = ProducerResponsibility
		prd.backup.Send(logmsg)
		return
	}
	if err := prd.mp.Produce(logmsg); err != nil {
		if sendToBackup {
			logmsg["from"] = ProducerResponsibility
			prd.backup.Send(logmsg)
		}
	}
	return
}

func RegisterMessageProducerFactory(fact func() sputnik.MessageProducer) {
	mpf = fact
}

var mpf func() sputnik.MessageProducer
