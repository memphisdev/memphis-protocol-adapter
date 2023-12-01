package main

import (
	"embed"

	"github.com/g41797/sputnik/sidecar"
	"github.com/memphisdev/memphis-protocol-adapter/pkg/adapter"

	// Attach blocks and plugins to the process:
	_ "github.com/memphisdev/memphis-protocol-adapter/pkg/grpcadapter"
	// ........................................
)

//go:embed conf
var embconf embed.FS

func main() {
	cleanUp, _ := sidecar.UseEmbeddedConfiguration(&embconf)
	defer cleanUp()

	sidecar.Start(new(adapter.BrokerConnector))
}
