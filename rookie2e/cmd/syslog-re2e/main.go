package main

import (
	"github.com/g41797/sputnik/sidecar"
	"github.com/memphisdev/memphis-protocol-adapter/pkg/adapter"

	// Attach blocks packages to the process:
	_ "github.com/memphisdev/memphis-protocol-adapter/rookie2e/pkg/syslogre2e"
)

func main() {
	sidecar.Start(new(adapter.BrokerConnector))
}
