package syslogre2e

import (
	"context"
	"fmt"
	"time"

	"github.com/g41797/sputnik"
	"github.com/memphisdev/memphis-protocol-adapter/pkg/syslogblocks"
	"github.com/memphisdev/memphis.go"
)

const (
	SyslogConsumerName           = "syslogconsumer"
	SyslogConsumerResponsibility = "syslogconsumer"
)

func SyslogConsumerDescriptor() sputnik.BlockDescriptor {
	return sputnik.BlockDescriptor{Name: SyslogClientName, Responsibility: SyslogClientResponsibility}
}

func init() {
	sputnik.RegisterBlockFactory(SyslogConsumerName, syslogConsumerBlockFactory)
}

func syslogConsumerBlockFactory() *sputnik.Block {
	cons := new(consumer)
	block := sputnik.NewBlock(
		sputnik.WithInit(cons.init),
		sputnik.WithRun(cons.run),
		sputnik.WithFinish(cons.finish),
		sputnik.WithOnConnect(cons.brokerConnected),
		sputnik.WithOnDisconnect(cons.brokerDisconnected),
	)
	return block
}

type consumer struct {
	conf    syslogblocks.MsgPrdConfig
	sender  sputnik.BlockCommunicator
	started bool
	mconn   *memphis.Conn
	mst     *memphis.Station
	mcons   *memphis.Consumer

	stop chan struct{}
	done chan struct{}
	conn chan struct{}
	dscn chan struct{}
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

waitBroker:
	for {
		select {
		case <-cons.stop:
			return
		case <-cons.conn:
			break waitBroker
		}
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

waitConsumer:
	for {
		select {
		case <-cons.stop:
			return
		case <-cons.dscn:
			return
		case <-ticker.C:
			if cons.startConsume() {
				break waitConsumer
			}
		}
	}

	cons.runLoop()

	cons.stopConsume()
	return
}

func (cons *consumer) runLoop() {
	for {
		select {
		case <-cons.stop:
			return
		case <-cons.dscn:
			return

		}
	}

	return
}

func (cons *consumer) startConsume() bool {

	mc, err := memphis.Connect(cons.conf.MEMPHIS_HOST, cons.conf.MEMPHIS_USER, memphis.Password(cons.conf.MEMPHIS_PSWRD))

	if err != nil {
		return false
	}

	st, err := mc.CreateStation(cons.conf.STATION)
	if err != nil {
		mc.Close()
		return false
	}

	mconsumer, err := st.CreateConsumer(SyslogConsumerResponsibility)
	if err != nil {
		mc.Close()
		return false
	}

	cons.mconn = mc
	cons.mst = st
	cons.mcons = mconsumer

	cons.mcons.Consume(cons.processMessages)
	cons.startLog()
	cons.started = true
	return true
}

func (cons *consumer) stopConsume() {
	if !cons.started {
		return
	}
	cons.mcons.StopConsume()
	cons.mconn.Close()
	cons.stopLog()
	return
}

func (cons *consumer) startLog() {
	msg := sputnik.Msg{}
	msg["name"] = "start"
	cons.sender.Send(msg)
	return
}

func (cons *consumer) stopLog() {
	msg := sputnik.Msg{}
	msg["name"] = "stop"
	cons.sender.Send(msg)
	return
}

func (cons *consumer) processMessages(msgs []*memphis.Msg, err error, ctx context.Context) {
	if err != nil {
		return
	}

	fmt.Printf("Consumed %d messages \n", len(msgs))

	for _, msg := range msgs {

		data := string(msg.Data())
		headers := msg.GetHeaders()
		if len(headers) == 0 {
			continue
		}

		smsg := sputnik.Msg{}
		smsg["name"] = "consumed"
		smsg["consumed"] = headers
		smsg["data"] = data
		cons.sender.Send(smsg)
	}
	ackAll(msgs)

	// ??? cons.mcons.Consume(cons.processMessages) ???
}

func ackAll(msgs []*memphis.Msg) {
	for _, msg := range msgs {
		msg.Ack()
	}
}
