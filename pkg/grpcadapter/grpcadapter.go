package memphisgrpc

import (
	"fmt"

	"github.com/g41797/memphisgrpc"
	"github.com/g41797/sputnik"
)

const (
	grpcadapterName           = "grpcadapter"
	grpcadapterResponsibility = "grpcadapter"
)

func init() {
	sputnik.RegisterBlockFactory(grpcadapterName, grpcadapterBlockFactory)
}

func grpcadapterBlockFactory() *sputnik.Block {
	gr := new(grpcadapter)

	block := sputnik.NewBlock(
		sputnik.WithInit(gr.init),
		sputnik.WithRun(gr.run),
		sputnik.WithFinish(gr.finish),
		sputnik.WithOnConnect(gr.brokerConnected),
		sputnik.WithOnDisconnect(gr.brokerDisconnected),
	)
	return block
}

type grpcadapter struct {
	connected bool
	cfact     sputnik.ConfFactory
	rnr       *memphisgrpc.Runner
	bc        sputnik.BlockCommunicator
	conn      chan sputnik.ServerConnection
	stop      chan struct{}
	done      chan struct{}
	dscn      chan struct{}
	stopf     func()
}

// Init
func (gr *grpcadapter) init(fact sputnik.ConfFactory) error {
	gr.cfact = fact
	gr.rnr = new(memphisgrpc.Runner)
	if err := gr.rnr.Init(gr.cfact); err != nil {
		return err
	}

	gr.stop = make(chan struct{}, 1)
	gr.done = make(chan struct{}, 1)
	gr.conn = make(chan sputnik.ServerConnection, 1)
	gr.dscn = make(chan struct{}, 1)

	return nil
}

// Run
func (gr *grpcadapter) run(bc sputnik.BlockCommunicator) {

	gr.bc = bc
	defer close(gr.done)
	defer gr.stopRunner()

	for {
		select {
		case sc := <-gr.conn:
			gr.onConnect(sc)

		case <-gr.dscn:
			if err := gr.onDisconnect(); err != nil {
				return
			}
		case <-gr.stop:
			return
		}
	}
}

// Finish:
func (gr *grpcadapter) finish(init bool) {
	if init {
		return
	}

	close(gr.stop) // Cancel Run

	<-gr.done // Wait finish of Run
	return
}

// OnServerConnect:
func (gr *grpcadapter) brokerConnected(srvc sputnik.ServerConnection) {
	gr.conn <- srvc
	return
}

// OnServerDisconnect:
func (gr *grpcadapter) brokerDisconnected() {
	gr.dscn <- struct{}{}
	return
}

func (gr *grpcadapter) onConnect(sc sputnik.ServerConnection) {
	gr.stopRunner()

	stop, err := gr.rnr.Run()
	if err != nil {
		return
	}

	gr.stopf = stop
	gr.connected = true

}

func (gr *grpcadapter) onDisconnect() error {
	return fmt.Errorf("reconnect is not supported")
}

func (gr *grpcadapter) stopRunner() {
	if gr == nil {
		return
	}
	if gr.stopf != nil {
		gr.stopf()
		gr.stopf = nil
		gr.connected = false
	}
}
