package main

import (
	"github.com/memphisdev/memphis-protocol-adapter/pkg/adapter"

	// Attach blocks packages to the process:
	_ "github.com/memphisdev/memphis-protocol-adapter/pkg/syslogblocks"
)

func main() {
	adapter.Run()
}
