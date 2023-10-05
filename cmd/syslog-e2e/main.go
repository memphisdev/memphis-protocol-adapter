package main

import (
	"fmt"

	"github.com/g41797/sputnik/sidecar"
	"github.com/memphisdev/memphis-protocol-adapter/pkg/adapter"

	// Attach blocks and plugins to the process:
	_ "github.com/g41797/syslogsidecar"
	e2e "github.com/g41797/syslogsidecar/e2e"
	_ "github.com/memphisdev/memphis-protocol-adapter/e2e/syslog"
	_ "github.com/memphisdev/memphis-protocol-adapter/pkg/syslog"
)

func main() {

	stop, err := e2e.StartServices()

	if err != nil {
		fmt.Println(err)
		return
	}

	defer stop()

	sidecar.Start(new(adapter.BrokerConnector))
}
