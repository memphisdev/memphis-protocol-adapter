package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/memphisdev/memphis-protocol-adapter/pkg/adapter"
)

func main() {

	var confFolder string
	flag.StringVar(&confFolder, "cf", "", "Path of folder with config files")
	flag.Parse()

	_, err := adapter.StartRunner(confFolder)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return
}
