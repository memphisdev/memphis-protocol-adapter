package main

import (
	"github.com/g41797/sputnik/sidecar"
	"github.com/memphisdev/memphis-protocol-adapter/pkg/adapter"

	// Attach blocks and plugins to the process:
	_ "github.com/g41797/syslogsidecar"
	_ "github.com/g41797/syslogsidecar/e2e"
	_ "github.com/memphisdev/memphis-protocol-adapter/e2e/syslog"
	_ "github.com/memphisdev/memphis-protocol-adapter/pkg/syslog"
)

func main() {
	sidecar.Start(new(adapter.BrokerConnector))
}
