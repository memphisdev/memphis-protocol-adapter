package syslogre2e

import (
	"github.com/g41797/sputnik"
	"github.com/memphisdev/memphis-protocol-adapter/pkg/syslogblocks"
)

const (
	SyslogConsumerName           = "syslogconsumer"
	SyslogConsumerResponsibility = "syslogconsumer"
)

func SyslogConsumerDescriptor() sputnik.BlockDescriptor {
	return sputnik.BlockDescriptor{Name: SyslogClientName, Responsibility: SyslogClientResponsibility}
}

func init() {
	sputnik.RegisterBlockFactory(SyslogClientName, syslogClientBlockFactory)
}

func syslogConsumerBlockFactory() *sputnik.Block {
	return nil
}

type consumer struct {
	conf    syslogblocks.MsgPrdConfig
	sender  sputnik.BlockCommunicator
	started bool
	stop    chan struct{}
	done    chan struct{}
	conn    chan struct{}
	dscn    chan struct{}
}

// Init
func (cons *consumer) init(fact sputnik.ConfFactory) error {
	if err := fact(syslogblocks.MsgProducerConfigName, &cons.conf); err != nil {
		return err
	}

	cons.stop = make(chan struct{}, 1)
	cons.done = make(chan struct{}, 1)
	cons.conn = make(chan struct{}, 1)
	cons.dscn = make(chan struct{}, 1)

	return nil
}

// Finish:
func (cons *consumer) finish(init bool) {
	if init {
		return
	}

	close(cons.stop) // Cancel Run

	<-cons.done // Wait finish of Run
	return
}

// OnServerConnect:
func (cons *consumer) brokerConnected(_ sputnik.ServerConnection) {
	cons.conn <- struct{}{}
	return
}

// OnServerDisconnect:
func (cons *consumer) brokerDisconnected() {
	cons.dscn <- struct{}{}
	return
}

// Run
func (cons *consumer) run(bc sputnik.BlockCommunicator) {

	cons.sender, _ = bc.Communicator(SyslogClientResponsibility)

	defer close(cons.done)

	cons.stopConsume()
	return
}

func (cons *consumer) startConsume() bool {
	return false
}

func (cons *consumer) stopConsume() {
	return
}
