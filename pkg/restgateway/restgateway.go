package restgateway

import (
	"fmt"

	"github.com/g41797/sputnik"
	"github.com/memphisdev/memphis-rest-gateway/gateway"
	lgr "github.com/memphisdev/memphis-rest-gateway/logger"
	"github.com/nats-io/nats.go"
)

const (
	restgatewayName           = "restgateway"
	restgatewayResponsibility = "restgateway"
)

func syslogConsumerDescriptor() sputnik.BlockDescriptor {
	return sputnik.BlockDescriptor{Name: restgatewayName, Responsibility: restgatewayResponsibility}
}

func init() {
	sputnik.RegisterBlockFactory(restgatewayName, restgatewayBlockFactory)
}

func restgatewayBlockFactory() *sputnik.Block {
	rg := new(restgateway)

	block := sputnik.NewBlock(
		sputnik.WithInit(rg.init),
		sputnik.WithRun(rg.run),
		sputnik.WithFinish(rg.finish),
		sputnik.WithOnConnect(rg.brokerConnected),
		sputnik.WithOnDisconnect(rg.brokerDisconnected),
	)
	return block
}

type restgateway struct {
	connected bool
	cfact     sputnik.ConfFactory
	bc        sputnik.BlockCommunicator
	conn      chan sputnik.ServerConnection
	stop      chan struct{}
	done      chan struct{}
	dscn      chan struct{}
	stopf     func() error
}

// Init
func (rg *restgateway) init(fact sputnik.ConfFactory) error {
	rg.cfact = fact
	rg.stop = make(chan struct{}, 1)
	rg.done = make(chan struct{}, 1)
	rg.conn = make(chan sputnik.ServerConnection, 1)
	rg.dscn = make(chan struct{}, 1)

	return nil
}

// Finish:
func (rg *restgateway) finish(init bool) {
	if init {
		return
	}

	close(rg.stop) // Cancel Run

	<-rg.done // Wait finish of Run
	return
}

// OnServerConnect:
func (rg *restgateway) brokerConnected(srvc sputnik.ServerConnection) {
	rg.conn <- srvc
	return
}

// OnServerDisconnect:
func (rg *restgateway) brokerDisconnected() {
	rg.dscn <- struct{}{}
	return
}

// Run
func (rg *restgateway) run(bc sputnik.BlockCommunicator) {

	rg.bc = bc
	defer close(rg.done)
	defer rg.stopRunner()

	for {
		select {
		case sc := <-rg.conn:
			rg.onConnect(sc)

		case <-rg.dscn:
			if err := rg.onDisconnect(); err != nil {
				return
			}
		case <-rg.stop:
			return
		}

	}
}

func (rg *restgateway) onConnect(sc sputnik.ServerConnection) {
	rg.stopRunner()

	logr, conn := sc.(func() (*lgr.Logger, *nats.Conn))()
	stop, err := gateway.Run(rg.cfact, logr, conn)
	if err != nil {
		return
	}
	rg.stopf = stop
	rg.connected = true

}

func (rg *restgateway) onDisconnect() error {
	return fmt.Errorf("reconnect is not supported")
}

func (rg *restgateway) stopRunner() {
	if rg == nil {
		return
	}
	if rg.stopf != nil {
		rg.stopf()
		rg.stopf = nil
		rg.connected = false
	}
}
