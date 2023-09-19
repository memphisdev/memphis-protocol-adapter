package main

import (
	"github.com/g41797/sputnik"
	"github.com/memphisdev/memphis-protocol-adapter/pkg/adapter"

	// Attach blocks packages to the process:
	_ "github.com/memphisdev/memphis-protocol-adapter/rookie2e/pkg/syslogre2e"
)

func main() {
	sputnik.Start(new(adapter.BrokerConnector))
}
