package syslogblocks

import "github.com/g41797/sputnik"

const MsgProducerConfigName = ProducerName

type MsgPrdConfig struct {
	PRODUCER             string
	STATION              string
	MEMPHIS_HOST         string
	USER_PASS_BASED_AUTH bool
	ROOT_USER            string
	ROOT_PASSWORD        string
	CONNECTION_TOKEN     string
	CLIENT_CERT_PATH     string
	CLIENT_KEY_PATH      string
	ROOT_CA_PATH         string
	CLOUD_ENV            bool
	DEBUG                bool
	DEV_ENV              bool
}

func newMsgProducer() MsgProducer {
	return &msgProducer{}
}

var _ MsgProducer = &msgProducer{}

type msgProducer struct {
}

func (mpr *msgProducer) Connect(cf sputnik.ConfFactory) error {
	return nil
}

func (mpr *msgProducer) Disconnect() {

}

func (mpr *msgProducer) Produce(msg sputnik.Msg) error {
	return nil
}
