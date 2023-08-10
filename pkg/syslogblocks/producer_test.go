package syslogblocks

import (
	"fmt"

	"github.com/g41797/kissngoqueue"
	"github.com/g41797/sputnik"
)

var _ MsgProducer = &MockMsgProducer{}

type MockMsgProducer struct {
	q *kissngoqueue.Queue[sputnik.Msg]
}

func newMMP(q *kissngoqueue.Queue[sputnik.Msg]) MsgProducer {
	res := MockMsgProducer{}
	res.q = q
	return &res
}

func (mp *MockMsgProducer) Connect(sputnik.ConfFactory) error {
	return nil
}

func (mp *MockMsgProducer) Disconnect() {
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
