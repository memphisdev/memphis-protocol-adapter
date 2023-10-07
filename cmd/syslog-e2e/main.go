package main

import (
	"fmt"

	// sputnik framework:
	"github.com/g41797/sputnik/sidecar"
	"github.com/g41797/starter"

	// memphis connector plugin
	"github.com/memphisdev/memphis-protocol-adapter/pkg/adapter"

	// Attach blocks and plugins to the process via blank imports:

	// 		syslog-adapter:
	// 			syslogsidecar blocks: receiver|producer|client|consumer
	_ "github.com/g41797/syslogsidecar"

	// 			memphis syslog plugins
	// 				msgconsumer
	_ "github.com/memphisdev/memphis-protocol-adapter/e2e/syslog"
	//				msgproducer
	_ "github.com/memphisdev/memphis-protocol-adapter/pkg/syslog"
)

func main() {

	stop, err := starter.StartServices()

	if err != nil {
		fmt.Println(err)
		return
	}

	defer stop()

	sidecar.Start(new(adapter.BrokerConnector))
}
