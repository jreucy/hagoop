package main

import (
	"strconv"
	"mrlib"
	"net"
	"fmt"
)

const (
	verbosity = mrlib.Verbosity
	MIN_JOB_SIZE = 1000  // change
	MAX_JOB_SIZE = 10000 // change
)

func main() {

	// "./server port"
	if len(os.Args) != 2 { return }

	port := os.Args[1]
	_ , err := strconv.Atoi(os.Args[1])
	if err != nil { /* do something */ }
	laddr := ":" + port

	// parse directory and save file name, starting/ending line numbers

	// connect to server with TCP
	serverAddr, err := net.ResolveTCPAddr("tcp", laddr)
	if err != nil { /* do something */ }
	serverListener, err := net.ListenTCP("tcp", serverAddr) // maybe change nil to something
	if err != nil { /* do something */ }

	// assuming a single request
	// later put in own goroutine
	conn, err := serverListener.AcceptTCP()
	if err != nil { /* do something */ }

	for {
		// Read in incoming requests from both request and worker clients

		// schedule and assign pending jobs to available workers

	}
}