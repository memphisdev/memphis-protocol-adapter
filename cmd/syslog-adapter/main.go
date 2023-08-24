package main

import (
	"github.com/memphisdev/memphis-protocol-adapter/pkg/adapter"

	// Blank imports:
	// "...on occasion we must import a package merely for the side effects
	// of doing so: evaluation of the initializer expressions of its
	// package-level variables and execution of its init functions..."

	// Attach blocks packages to the process:
	_ "github.com/memphisdev/memphis-protocol-adapter/pkg/syslogblocks"
)

func main() {
	adapter.Run()
}
