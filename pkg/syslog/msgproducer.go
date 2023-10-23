package syslog

import (
	"fmt"
	"strings"

	"github.com/g41797/sputnik"
	"github.com/g41797/sputnik/sidecar"
	"github.com/g41797/syslogsidecar"
	"github.com/memphisdev/memphis.go"
)

func init() {
	syslogsidecar.RegisterMessageProducerFactory(newMsgProducer)
}

const MsgProducerConfigName = syslogsidecar.ProducerName

type MsgPrdConfig struct {
	MEMPHIS_HOST   string
	MEMPHIS_USER   string
	MEMPHIS_PSWRD  string
	PRODUCER       string
	STATION        string
	RETENTIONTYPE  string
	RETENTIONVALUE int
}

func newMsgProducer() sidecar.MessageProducer {
	return &msgProducer{}
}

type msgProducer struct {
	conf     MsgPrdConfig
	mc       *memphis.Conn
	producer *memphis.Producer
}

func (mpr *msgProducer) Connect(cf sputnik.ConfFactory, _ sputnik.ServerConnection) error {
	err := cf(MsgProducerConfigName, &mpr.conf)
	if err != nil {
		return err
	}

	mpr.mc, err = memphis.Connect(mpr.conf.MEMPHIS_HOST, mpr.conf.MEMPHIS_USER, memphis.Password(mpr.conf.MEMPHIS_PSWRD))

	if err != nil {
		mpr.mc = nil
		return err
	}

	_, err = CreateStation(mpr.mc, &mpr.conf)
	if err != nil {
		mpr.Disconnect()
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

	if mpr.producer != nil {
		mpr.producer.Destroy()
		mpr.producer = nil
	}

	mpr.mc.Close()
	mpr.mc = nil
	return
}

func (mpr *msgProducer) Produce(msg sputnik.Msg) error {

	defer syslogsidecar.Put(msg)

	if mpr.mc == nil {
		return fmt.Errorf("connection with broker does not exist")
	}

	if !mpr.mc.IsConnected() {
		return fmt.Errorf("does not connected with broker")
	}

	hdrs := memphis.Headers{}
	hdrs.New()

	putToheader := func(name string, value string) error {
		return hdrs.Add(name, value)
	}

	if err := syslogsidecar.Unpack(msg, putToheader); err != nil {
		return err
	}

	err := mpr.producer.Produce("", memphis.MsgHeaders(hdrs))

	return err
}

func CreateStation(mc *memphis.Conn, conf *MsgPrdConfig) (*memphis.Station, error) {
	rt, rv := retentionParams(conf)
	st, err := mc.CreateStation(conf.STATION, memphis.RetentionTypeOpt(rt), memphis.RetentionVal(rv))

	if err != nil {
		return nil, err
	}

	return st, nil
}

var retentiontypes = []string{"MaxMessageAgeSeconds", "Messages", "Bytes", "AckBased"}

func retentionParams(conf *MsgPrdConfig) (rt memphis.RetentionType, rv int) {
	defaultOpts := memphis.GetStationDefaultOptions()

	rt = defaultOpts.RetentionType
	rv = conf.RETENTIONVALUE
	if rv == 0 {
		rv = defaultOpts.RetentionVal
	}

	for i, val := range retentiontypes {
		if strings.ToUpper(val) == strings.ToUpper(conf.RETENTIONTYPE) {
			rt = memphis.RetentionType(i)
			return rt, rv
		}
	}

	return rt, rv
}
