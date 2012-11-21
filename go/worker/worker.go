package main

import (
	"strconv"
	"mrlib"
	"net"
	"fmt"
)

const (
	verbosity = mrlib.Verbosity
)

func main() {

	// "./worker host:port"
	if len(os.Args) != 2 { return }

	hostport := os.Args[1]

	// Connect to server using TCP

	// Let server know worker is available

	for {

		// Read in Map or Reduce requests from server

		// perform job

		// send results back to server

	}

}