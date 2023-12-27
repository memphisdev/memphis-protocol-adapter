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
	conf          MsgPrdConfig
	rt            memphis.RetentionType
	rv            int
	mc            *memphis.Conn
	pr4stat       map[string]*memphis.Producer
	usesyslogconf bool
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

	mpr.rt, mpr.rv = retentionParams(&mpr.conf)

	mpr.pr4stat = make(map[string]*memphis.Producer)

	err = mpr.CreateProducerAndStation(mpr.conf.STATION)

	if err != nil {
		mpr.Disconnect()
		return err
	}

	if _, err = syslogsidecar.AllTargets(); err == nil {
		mpr.usesyslogconf = true
	}

	return nil
}

func (mpr *msgProducer) CreateProducerAndStation(station string) error {

	if _, exists := mpr.pr4stat[station]; exists {
		return nil
	}

	st, err := mpr.mc.CreateStation(station, memphis.RetentionTypeOpt(mpr.rt), memphis.RetentionVal(mpr.rv))

	if err != nil {
		return err
	}

	p, err := st.CreateProducer(mpr.conf.PRODUCER)
	if err != nil {
		return err
	}

	mpr.pr4stat[station] = p

	return nil
}

func (mpr *msgProducer) Disconnect() {
	if mpr.mc == nil {
		return
	}

	for _, producer := range mpr.pr4stat {
		if producer != nil {
			producer.Destroy()
		}
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

	if !mpr.usesyslogconf {
		return mpr.pr4stat[mpr.conf.STATION].Produce("", memphis.MsgHeaders(hdrs))
	}

	stations, _ := syslogsidecar.Targets(msg)

	if len(stations) == 0 {
		return mpr.pr4stat[mpr.conf.STATION].Produce("", memphis.MsgHeaders(hdrs))
	}

	for _, station := range stations {
		if err := mpr.CreateProducerAndStation(station); err != nil {
			return err
		}
		if err := mpr.pr4stat[station].Produce("", memphis.MsgHeaders(hdrs)); err != nil {
			return err
		}
	}

	return nil
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
