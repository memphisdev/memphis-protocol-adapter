package syslogblocks

import (
	"fmt"

	"github.com/g41797/sputnik"
	"github.com/memphisdev/memphis.go"
)

const MsgProducerConfigName = ProducerName

type MsgPrdConfig struct {
	MEMPHIS_HOST  string
	MEMPHIS_USER  string
	MEMPHIS_PSWRD string
	PRODUCER      string
	STATION       string
}

func newMsgProducer() MsgProducer {
	return &msgProducer{}
}

var _ MsgProducer = &msgProducer{}

type msgProducer struct {
	conf     MsgPrdConfig
	mc       *memphis.Conn
	producer *memphis.Producer
}

func (mpr *msgProducer) Connect(cf sputnik.ConfFactory) error {
	err := cf(MsgProducerConfigName, &mpr.conf)
	if err != nil {
		return err
	}

	mpr.mc, err = memphis.Connect(mpr.conf.MEMPHIS_HOST, mpr.conf.MEMPHIS_USER, memphis.Password(mpr.conf.MEMPHIS_PSWRD))

	if err != nil {
		mpr.mc = nil
		return err
	}

	p, err := mpr.mc.CreateProducer(mpr.conf.STATION, mpr.conf.PRODUCER)

	if err != nil {
		mpr.Disconnect()
		return err
	}

	mpr.producer = p
	return nil
}

func (mpr *msgProducer) Disconnect() {
	if mpr.mc == nil {
		return
	}

	mpr.mc.Close()
	mpr.mc = nil
	return
}

func (mpr *msgProducer) Produce(msg sputnik.Msg) error {
	if mpr.mc == nil {
		return fmt.Errorf("connection with broker does not exist")
	}

	if !mpr.mc.IsConnected() {
		return fmt.Errorf("does not connected with broker")
	}

	hdrs := memphis.Headers{}
	hdrs.New()

	for k, v := range msg {
		vstr, ok := v.(string)
		if !ok {
			continue
		}
		if err := hdrs.Add(k, vstr); err != nil {
			return err
		}
	}

	err := mpr.producer.Produce("", memphis.MsgHeaders(hdrs))

	return err
}
