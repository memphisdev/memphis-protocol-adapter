package sysloge2e

import (
	"context"
	"time"

	"github.com/g41797/sputnik"
	"github.com/g41797/sputnik/sidecar"
	"github.com/g41797/syslogsidecar/e2e"
	"github.com/memphisdev/memphis-protocol-adapter/pkg/adapter"
	"github.com/memphisdev/memphis-protocol-adapter/pkg/syslog"
	"github.com/memphisdev/memphis.go"
)

func init() {
	e2e.RegisterMessageConsumerFactory(newMsgConsumer)
}

const MsgConsumerConfigName = syslog.MsgProducerConfigName

type msgConsumer struct {
	conf    syslog.MsgPrdConfig
	mconn   *memphis.Conn
	mst     *memphis.Station
	mcons   *memphis.Consumer
	sender  sputnik.BlockCommunicator
	started bool
}

func newMsgConsumer() sidecar.MessageConsumer {
	return new(msgConsumer)
}

func (mcn *msgConsumer) Connect(cf sputnik.ConfFactory, shared sputnik.ServerConnection) error {

	if err := cf(MsgConsumerConfigName, &mcn.conf); err != nil {
		return err
	}

	if err := mcn.prepare(); err != nil {
		return err
	}

	lgrf, ok := shared.(func() *adapter.Logger)
	if ok {
		lgrf().Noticef("Syslog consumer started")
	}

	return nil
}

func (cons *msgConsumer) Consume(sender sputnik.BlockCommunicator) error {
	if cons.started {
		return nil
	}

	cons.sender = sender
	cons.mcons.Consume(cons.processMessages)
	cons.startTest()
	cons.started = true
	return nil
}

func (cons *msgConsumer) Disconnect() {
	if cons == nil {
		return
	}

	if !cons.started {
		return
	}
	if cons.mcons != nil {
		cons.mcons.StopConsume()
		cons.mconn.Close()
		cons.mcons = nil
		cons.mconn = nil
	}
	cons.stopTest()
	cons.started = false
	return
}

func (cons *msgConsumer) prepare() error {

	mc, err := memphis.Connect(cons.conf.MEMPHIS_HOST, cons.conf.MEMPHIS_USER, memphis.Password(cons.conf.MEMPHIS_PSWRD))

	if err != nil {
		return err
	}

	st, err := mc.CreateStation(cons.conf.STATION)
	if err != nil {
		mc.Close()
		return err
	}

	st.Destroy()

	st, _ = mc.CreateStation(cons.conf.STATION)

	mconsumer, err := st.CreateConsumer(e2e.SyslogConsumerResponsibility, memphis.PullInterval(50*time.Millisecond), memphis.BatchSize(1000), memphis.BatchMaxWaitTime(time.Second))
	if err != nil {
		mc.Close()
		return err
	}

	cons.mconn = mc
	cons.mst = st
	cons.mcons = mconsumer
	return nil
}

func (cons *msgConsumer) processMessages(msgs []*memphis.Msg, err error, ctx context.Context) {
	if err != nil {
		return
	}

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
}

func (cons *msgConsumer) startTest() {
	msg := sputnik.Msg{}
	msg["name"] = "start"
	cons.sender.Send(msg)
}

func (cons *msgConsumer) stopTest() {
	msg := sputnik.Msg{}
	msg["name"] = "stop"
	cons.sender.Send(msg)
}

func ackAll(msgs []*memphis.Msg) {
	for _, msg := range msgs {
		msg.Ack()
	}
}
