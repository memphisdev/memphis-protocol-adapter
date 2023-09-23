package main

import (
	"github.com/g41797/sputnik/sidecar"
	"github.com/memphisdev/memphis-protocol-adapter/pkg/adapter"

	// Attach blocks to the process:
	_ "github.com/g41797/syslogsidecar"
	_ "github.com/g41797/syslogsidecar/e2e"
)

func main() {
	sidecar.Start(new(adapter.BrokerConnector))
}
