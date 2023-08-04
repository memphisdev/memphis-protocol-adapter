package syslogblocks_test

import (
	"testing"

	"github.com/g41797/kissngoqueue"
	"github.com/g41797/sputnik"
	"github.com/memphisdev/memphis-protocol-adapter/pkg/syslogblocks"
)

func Test_Init_Finish(t *testing.T) {

	q := kissngoqueue.NewQueue[sputnik.Msg]()
	mc := newCommunicator(q)

	srv := syslogblocks.NewServer(defaultServerConfiguration(), mc)
	defer stopServer(t, srv)
	err := srv.Init()
	if err != nil {
		t.Errorf("Init error %v", err)
	}
}

func stopServer(t *testing.T, srv *syslogblocks.Server) {
	err := srv.Finish()
	if err != nil {
		t.Errorf("stop server error %v", err)
	}
}

func defaultServerConfiguration() syslogblocks.SyslogConfiguration {
	result := syslogblocks.SyslogConfiguration{}
	result.ADDRTCP = "0.0.0.0:5141"
	result.SEVERITYLEVEL = 7
	return result
}

var _ sputnik.BlockCommunicator = &MockCommunicator{}

type MockCommunicator struct {
	q *kissngoqueue.Queue[sputnik.Msg]
}

func newCommunicator(q *kissngoqueue.Queue[sputnik.Msg]) *MockCommunicator {
	mc := new(MockCommunicator)
	mc.q = q
	return mc
}

func (mc *MockCommunicator) Communicator(resp string) (bc sputnik.BlockCommunicator, exists bool) {
	return nil, false
}

func (mc *MockCommunicator) Descriptor() sputnik.BlockDescriptor {
	return sputnik.BlockDescriptor{}
}

func (mc *MockCommunicator) Send(msg sputnik.Msg) bool {
	if msg == nil {
		return false
	}

	if mc.q == nil {
		return false
	}

	sok := mc.q.PutMT(msg)

	return sok
}
