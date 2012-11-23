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
	serverAddr, err := net.ResolveTCPAddr("tcp", hostport)
	if err != nil { /* do something */ }
	conn, err := net.DialTCP("tcp", nil, serverAddr) // maybe change nil to something
	if err != nil { /* do something */ }

	// Let server know worker is available

	for {

		// Read in Map or Reduce requests from server

		// perform job

		// send results back to server

	}

}