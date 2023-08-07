package syslogblocks

import (
	"strconv"
	"testing"

	syslogclient "github.com/RackSec/srslog"
	"github.com/g41797/kissngoqueue"
	"github.com/g41797/sputnik"
)

func defaultServerConfiguration() SyslogConfiguration {
	result := SyslogConfiguration{}
	result.ADDRTCP = "127.0.0.1:5141"
	result.SEVERITYLEVEL = 7
	return result
}

var seqNumber int

func next() int {
	seqNumber += 1
	return seqNumber
}

type srvTest struct {
	t          *testing.T
	q          *kissngoqueue.Queue[sputnik.Msg]
	syslogconf SyslogConfiguration
	srv        *Server
	client     *client
}

func newTest(t *testing.T) *srvTest {
	result := new(srvTest)
	result.q = kissngoqueue.NewQueue[sputnik.Msg]()
	return result
}

func (test *srvTest) start() {
	test.syslogconf = defaultServerConfiguration()
	test.srv = NewServer(test.syslogconf)

	err := test.srv.Init()
	if err != nil {
		test.t.Errorf("Init error %v", err)
	}

	test.client = newClient()

	err = test.client.init()
	if err != nil {
		test.t.Errorf("Init client error %v", err)
	}
	test.srv.SetupHandling(newCommunicator(test.q))

	err = test.srv.Start()
	if err != nil {
		test.t.Errorf("Start syslogd error %v", err)
	}
}

func (test *srvTest) stop() {
	test.client.finish()
	err := test.srv.Finish()
	if err != nil {
		test.t.Errorf("stop server error %v", err)
	}
}

func (test *srvTest) exchange() {
	next := next()
	logText := strconv.Itoa(next)
	err := test.client.log(logText)
	if err != nil {
		test.t.Errorf("send log message error %v", err)
	}

	msg, ok := test.q.Get()

	if !ok {
		test.t.Errorf("failed receive from test queue")
	}

	recvlog, _ := msg["message"].(string)

	if logText != recvlog {
		test.t.Errorf("Expected %s Received %s", logText, recvlog)
	}
}

func (test *srvTest) stopServer() {
	err := test.srv.Finish()
	if err != nil {
		test.t.Errorf("stop server error %v", err)
	}
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

type client struct {
	syslogconf SyslogConfiguration
	w          *syslogclient.Writer
}

func newClient() *client {
	result := new(client)
	result.syslogconf = defaultServerConfiguration()
	return result
}

func (cl *client) init() error {
	w, err := syslogclient.Dial("tcp", cl.syslogconf.ADDRTCP, syslogclient.LOG_ALERT, "test-client")
	if err != nil {
		return err
	}
	w.SetFormatter(syslogclient.RFC5424Formatter)
	cl.w = w
	return nil
}

func (cl *client) log(msg string) error {
	return cl.w.Alert(msg)
}

func (cl *client) finish() {
	if cl == nil {
		return
	}
	if cl.w == nil {
		return
	}

	cl.w.Close()
	return
}

func Test_Init_Finish(t *testing.T) {
	test := newTest(t)
	test.start()
	test.stop()
}

func Test_Exchange(t *testing.T) {
	test := newTest(t)
	test.start()

	defer test.stop()

	for i := 0; i < 1000; i++ {
		test.exchange()
	}
}
