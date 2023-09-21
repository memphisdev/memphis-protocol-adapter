package main

import (
	"github.com/g41797/sputnik/sidecar"
	"github.com/memphisdev/memphis-protocol-adapter/pkg/adapter"

	// Attach blocks packages to the process:
	_ "github.com/g41797/syslogsidecar"
)

func main() {
	sidecar.Start(new(adapter.BrokerConnector))
}
