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

	// "./server port"
	if len(os.Args) != 2 { return }

	
	for {
		// Read in incoming requests from both request and worker clients

		// schedule and assign pending jobs to available workers

	}
}