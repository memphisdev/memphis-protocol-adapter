package adapter

import (
	"fmt"
	"os"
	"time"

	"github.com/g41797/sputnik"
)

const brokerCheckTimeOut = time.Second

type Runner struct {
	// ShootDown
	kill sputnik.ShootDown
	// Signalling channel
	done chan struct{}
}

func StartRunner(confFolder string) (*Runner, error) {
	info, err := prepare(confFolder)
	if err != nil {
		return nil, err
	}
	rnr := new(Runner)
	err = rnr.Start(info)
	if err != nil {
		return nil, err
	}
	return rnr, nil
}

func (rnr *Runner) Stop() {
	if rnr == nil {
		return
	}
	if rnr.kill == nil {
		return
	}

	select {
	case <-rnr.done:
		return
	default:
	}

	rnr.kill()

	return
}

type runnerInfo struct {
	cfact     sputnik.ConfFactory
	cnt       sputnik.ServerConnector
	appBlocks []sputnik.BlockDescriptor
}

func prepare(confFolder string) (*runnerInfo, error) {
	info, err := os.Stat(confFolder)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not the folder", confFolder)
	}

	ri := runnerInfo{}

	ri.cfact = ConfigFactory(confFolder)

	ri.cnt = new(BrokerConnector)

	ri.appBlocks, err = ReadAppBlocks(confFolder)

	return &ri, err
}

func (rnr *Runner) Start(ri *runnerInfo) error {

	sp, err := sputnik.NewSputnik(
		sputnik.WithAppBlocks(ri.appBlocks),
		sputnik.WithConfFactory(ri.cfact),
		sputnik.WithConnector(ri.cnt, brokerCheckTimeOut),
	)

	if err != nil {
		return err
	}

	launch, kill, err := sp.Prepare()

	if err != nil {
		return err
	}
	rnr.kill = kill
	rnr.done = make(chan struct{})

	go func(l sputnik.Launch, done chan struct{}) {
		l()
		close(done)
	}(launch, rnr.done)

	return nil
}

func ReadAppBlocks(confFolder string) ([]sputnik.BlockDescriptor, error) {
	return nil, nil
}
