package syslogblocks

import (
	"strconv"
	"testing"
	"time"

	"github.com/g41797/kissngoqueue"
	"github.com/g41797/sputnik"
	"github.com/memphisdev/memphis-protocol-adapter/pkg/adapter"
)

// Satellite has 2 app. blocks:
var blkList []sputnik.BlockDescriptor = []sputnik.BlockDescriptor{
	// memphis events Producer (simulated by echo block)
	{sputnik.EchoBlockName, ProducerResponsibility},
	// syslog Receiver
	{ReceiverName, ReceiverResponsibility},
}

type recvTest struct {
	// Config path & factory
	confFolderPath string
	cfact          sputnik.ConfFactory
	syslogconf     SyslogConfiguration
	// Test queue
	q *kissngoqueue.Queue[sputnik.Msg]
	// Launcher
	launch sputnik.Launch
	// ShootDown
	kill sputnik.ShootDown
	// Signalling channel
	done chan struct{}
	// ServerConnector
	conntr sputnik.DummyConnector
	to     time.Duration

	client *client
}

func newRecvTest() *recvTest {
	res := new(recvTest)
	res.confFolderPath = "./_conf_test/"
	res.cfact = adapter.ConfigFactory(res.confFolderPath)
	res.q = kissngoqueue.NewQueue[sputnik.Msg]()
	res.conntr = sputnik.DummyConnector{}
	res.to = time.Millisecond * 100
	return res
}

// Registration of factories for test environment
// For this case init() isn't used
// use this pattern for the case when you don't need
// dynamic registration: all blocks (and factories) are
// known in advance.
func (rt *recvTest) factories() sputnik.BlockFactories {
	res := make(sputnik.BlockFactories)

	finfct, _ := sputnik.Factory(sputnik.DefaultFinisherName)
	confct, _ := sputnik.Factory(sputnik.DefaultConnectorName)
	echoFact := sputnik.EchoBlockFactory(rt.q)

	factList := []struct {
		name string
		fact sputnik.BlockFactory
	}{
		{sputnik.DefaultFinisherName, finfct},
		{sputnik.DefaultConnectorName, confct},
		{sputnik.EchoBlockName, echoFact},
		{ReceiverName, receiverBlockFactory},
	}

	for _, fd := range factList {
		sputnik.RegisterBlockFactoryInner(fd.name, fd.fact, res)
	}
	return res
}

func recvSputnik(rt *recvTest) sputnik.Sputnik {
	sp, _ := sputnik.NewSputnik(
		sputnik.WithConfFactory(rt.cfact),
		sputnik.WithAppBlocks(blkList),
		sputnik.WithBlockFactories(rt.factories()),
		sputnik.WithConnector(&rt.conntr, rt.to),
	)
	return *sp
}

// Run Launcher on dedicated goroutine
// Test controls execution via sputnik API
// Results received using queue
func (rt *recvTest) run(t *testing.T) {

	if err := rt.cfact(ReceiverName, &rt.syslogconf); err != nil {
		t.Fatalf("failed get configuration")
	}

	rt.client = newClientWIthConfig(rt.syslogconf)

	rt.done = make(chan struct{})
	if rt.launch == nil {
		t.Fatalf("nil rt.launch")
	}

	go func(l sputnik.Launch, done chan struct{}, client *client) {
		if l == nil {
			panic("nil launcher")
		}
		l()

		if client != nil {
			client.finish()
		}

		close(done)
	}(rt.launch, rt.done, rt.client)

	return
}

func (rt *recvTest) exchange(t *testing.T) {
	next := next()
	logText := strconv.Itoa(next)
	err := rt.client.log(logText)
	if err != nil {
		t.Errorf("send log message error %v", err)
	}

	msg, ok := rt.q.Get()

	if !ok {
		t.Errorf("failed receive from test queue")
	}

	recvlog, _ := msg["message"].(string)

	if logText != recvlog {
		t.Errorf("Expected %s Received %s", logText, recvlog)
	}
}

func TestReceive_Prepare(t *testing.T) {

	rt := newRecvTest()

	sputnik := recvSputnik(rt)

	_, kill, err := sputnik.Prepare()

	if err != nil {
		t.Errorf("Prepare error %v", err)
	}

	time.Sleep(time.Second)

	kill()

	return
}

func TestReceive_StartStop(t *testing.T) {

	rt := newRecvTest()

	sputnik := recvSputnik(rt)

	launch, kill, err := sputnik.Prepare()

	if err != nil {
		t.Fatalf("prepare error %v", err)
	}

	if launch == nil {
		t.Fatalf("nil launch")
	}

	rt.launch = launch

	rt.kill = kill

	rt.run(t)

	time.Sleep(time.Millisecond * 100)
	rt.conntr.SetState(true)

	if err = rt.client.init(); err != nil {
		t.Fatalf("failure of syslogd client %v", err)
	}

	time.Sleep(time.Millisecond * 100)
	rt.kill()

	return
}

func TestReceive_Exchange(t *testing.T) {

	rt := newRecvTest()

	sputnik := recvSputnik(rt)

	launch, kill, err := sputnik.Prepare()

	if err != nil {
		t.Fatalf("prepare error %v", err)
	}

	if launch == nil {
		t.Fatalf("nil launch")
	}

	rt.launch = launch

	rt.kill = kill

	rt.run(t)

	time.Sleep(time.Millisecond * 100)
	rt.conntr.SetState(true)

	if err = rt.client.init(); err != nil {
		t.Fatalf("failure of syslogd client %v", err)
	}
	time.Sleep(time.Millisecond * 100)

	for i := 0; i < 50000; i++ {
		rt.exchange(t)
	}

	time.Sleep(time.Millisecond * 100)
	rt.kill()

	return
}
