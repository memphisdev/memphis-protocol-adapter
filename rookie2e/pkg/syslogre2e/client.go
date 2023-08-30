package syslogre2e

import (
	"fmt"
	"os"
	"strconv"

	"github.com/RackSec/srslog"
	"github.com/g41797/sputnik"
	"github.com/memphisdev/memphis-protocol-adapter/pkg/syslogblocks"
	"github.com/memphisdev/memphis.go"
)

func NewLogWriter(cnf syslogblocks.SyslogConfiguration, rfcForm srslog.Formatter) (*srslog.Writer, error) {
	w, err := srslog.Dial("tcp", cnf.ADDRTCP, srslog.LOG_ALERT, "re2e")
	if err != nil {
		return nil, err
	}
	w.SetFormatter(rfcForm)
	return w, nil
}

func SendQuit() error {
	cproc, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}
	err = cproc.Signal(os.Interrupt)
	return err
}

const MAX_LOG_MESSAGES = 1000

const (
	SENDOK = iota + 1
	CONSUMED
)

type client struct {
	conf    syslogblocks.SyslogConfiguration
	loggers []*srslog.Writer
	bc      sputnik.BlockCommunicator

	started  bool
	currIndx int
	states   []int
	successN int

	startFlow chan struct{}
	stopFlow  chan struct{}
	brokMsg   chan *memphis.Msg
	nextSend  chan struct{}

	stop chan struct{}
	done chan struct{}
}

// Init
func (cl *client) init(fact sputnik.ConfFactory) error {
	if err := fact(syslogblocks.ReceiverName, &cl.conf); err != nil {
		return err
	}

	cl.loggers = make([]*srslog.Writer, 2)
	cl.stop = make(chan struct{}, 1)
	cl.done = make(chan struct{}, 1)
	cl.startFlow = make(chan struct{}, 1)
	cl.stopFlow = make(chan struct{}, 1)
	cl.brokMsg = make(chan *memphis.Msg)
	cl.nextSend = make(chan struct{})

	return nil
}

// Finish:
func (cl *client) finish(init bool) {
	if init {
		return
	}

	close(cl.stop) // Cancel Run

	<-cl.done // Wait finish of Run
	return
}

// Run
func (cl *client) run(bc sputnik.BlockCommunicator) {

	cl.bc = bc

	defer close(cl.done)

loop:
	for {
		select {
		case <-cl.stop:
			break loop
		case <-cl.startFlow:
			cl.startflow()
		case <-cl.stopFlow:
			cl.stopflow()
		case brokermsg := <-cl.brokMsg:
			cl.update(brokermsg)
		case <-cl.nextSend:
			cl.sendNext()
		}
	}

	cl.closeLoggers()
	cl.report()

	return
}

// OnMsg:
func (cl *client) processBrokerMsg(brokermsg sputnik.Msg) {
	if brokermsg == nil {
		return
	}

	name, exists := brokermsg["name"]
	if !exists {
		return
	}

	switch name {
	case "start":
		cl.startFlow <- struct{}{}
	case "stop":
		cl.stopFlow <- struct{}{}
	case "consumed":
		brmsg, ok := brokermsg["consumed"].(*memphis.Msg)
		if ok && brmsg != nil {
			cl.brokMsg <- brmsg
		}
	}

	return
}

func (cl *client) startflow() {
	if cl.started {
		return
	}

	if err := cl.openLoggers(); err != nil {
		cl.stopflow()
		return
	}

	cl.started = true
	cl.states = make([]int, MAX_LOG_MESSAGES, MAX_LOG_MESSAGES)
	cl.sendNext()
}

func (cl *client) sendNext() {
	if !cl.started {
		return
	}

	if cl.currIndx >= MAX_LOG_MESSAGES {
		cl.stopflow()
		return
	}

	if err := cl.loggers[cl.currIndx%2].Warning(strconv.Itoa(cl.currIndx)); err != nil {
		cl.stopflow()
		return
	}

	cl.states[cl.currIndx] = SENDOK
	cl.currIndx++
	cl.nextSend <- struct{}{}

	return
}

func (cl *client) stopflow() {
	cl.started = false
	SendQuit()
}

func (cl *client) update(brmsg *memphis.Msg) {
	if brmsg == nil {
		return
	}

	hdrs := brmsg.GetHeaders()
	if hdrs == nil {
		return
	}

	rfc, ok := hdrs[syslogblocks.RFCFormatKey]
	if !ok {
		return
	}

	valName := "message"

	if rfc == syslogblocks.RFC3164 {
		valName = "content"
	}

	msgText, ok := hdrs[valName]
	if !ok {
		return
	}

	msgIndex, err := strconv.Atoi(msgText)
	if err != nil {
		return
	}

	if msgIndex >= MAX_LOG_MESSAGES {
		return
	}

	status := cl.states[msgIndex]
	if status == SENDOK {
		cl.states[msgIndex] = CONSUMED
		cl.successN++
	}
	return
}

func (cl *client) report() {
	fmt.Printf("\n\n\t\tWas send %d messages. Successfully consumed %d\n\n", cl.currIndx, cl.successN)
	return
}

func (cl *client) openLoggers() error {
	lgr, err := NewLogWriter(cl.conf, srslog.RFC3164Formatter)
	if err != nil {
		return err
	}
	cl.loggers[0] = lgr

	lgr, err = NewLogWriter(cl.conf, srslog.RFC5424Formatter)
	if err != nil {
		cl.loggers[0].Close()
		cl.loggers[0] = nil
		return err
	}
	cl.loggers[1] = lgr
	return nil
}

func (cl *client) closeLoggers() {
	for _, lgr := range cl.loggers {
		if lgr != nil {
			lgr.Close()
		}
	}
}
