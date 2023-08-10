package syslogblocks

import (
	"fmt"

	"github.com/g41797/kissngoqueue"
	"github.com/g41797/sputnik"
)

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

var _ MsgProducer = &MockMsgProducer{}

type MockMsgProducer struct {
	cn *sputnik.DummyConnector
	q  *kissngoqueue.Queue[sputnik.Msg]
}

func newMMP(cn *sputnik.DummyConnector, q *kissngoqueue.Queue[sputnik.Msg]) *MockMsgProducer {
	res := new(MockMsgProducer)
	res.cn = cn
	res.q = q
	return res
}

func (mp *MockMsgProducer) Connect(sputnik.ConfFactory) error {
	if mp.cn != nil {
		if mp.cn.IsConnected() {
			return nil
		}
	}
	return fmt.Errorf("not connected with broker")
}

func (mp *MockMsgProducer) Disconnect() {
	if mp.cn != nil {
		mp.cn.SetState(false)
	}
	return
}

func (mp *MockMsgProducer) Produce(msg sputnik.Msg) error {
	if mp.q == nil {
		return fmt.Errorf("q does not exist. wrong test environment")
	}
	if ok := mp.q.PutMT(msg); ok {
		return nil
	}

	return fmt.Errorf("q canceled")
}

func newMsgProducer() *MsgProducer {
	return nil
}
