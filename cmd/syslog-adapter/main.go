package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/memphisdev/memphis-protocol-adapter/pkg/adapter"

	// Blank imports:
	// "...on occasion we must import a package merely for the side effects
	// of doing so: evaluation of the initializer expressions of its
	// package-level variables and execution of its init functions..."

	// Attach syslogblocks package to the process:
	_ "github.com/memphisdev/memphis-protocol-adapter/pkg/syslogblocks"
)

func main() {

	var confFolder string
	flag.StringVar(&confFolder, "cf", "", "Path of folder with config files")
	flag.Parse()

	rnr, err := adapter.StartRunner(confFolder)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rnr.Wait()
	return
}

/*
func TryConnect() error {

	hostname := "localhost"
	username := "$memphis"
	creds := "memphis_memphis"

	var nc *nats.Conn
	var err error

	natsOpts := nats.Options{
		Url:            hostname + ":6666",
		AllowReconnect: true,
		MaxReconnect:   10,
		ReconnectWait:  3 * time.Second,
		Name:           "MEMPHIS HTTP LOGGER",
	}

	natsOpts.Password = creds
	natsOpts.User = username
	natsOpts.TLSConfig = nil

	nc, err = natsOpts.Connect()
	if err != nil {
		return err
	}

	nc.Close()

	return nil
}
*/
