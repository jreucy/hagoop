package main

import (
	"strconv"
	"mrlib"
	"net"
	"fmt"
)

func main() {

	// "./request host:port [file_name] [start_line] [end_line]"
	if len(os.Args) != 4 { return }

	hostport := os.Args[1]
	file_name := os.Args[2]
	start_line := strconv.Atoi(os.Args[3])
	end_line := strconv.Atoi(os.Args[4])

	// connect to server with TCP

	// create mapreduce request and write to server
	mr_request := []string{mrlib.MsgMAPREDUCE, file_name, start_line, end_line}
	mr_request_msg := strings.Join(mr_request, " ")

	// read answer from server and print to command-line
}