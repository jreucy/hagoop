package main

import (
	"strconv"
	"hagoop/go/mrlib"
	"net"
	"fmt"
)



func main() {

	// "./request host:port [file_directory] [file_name]"
	if len(os.Args) != 4 { return }


	hostport := os.Args[1]
	file_directory := os.Args[2]
	answer_file_name := os.Args[3]

	// connect to server with TCP
	serverAddr, err := net.ResolveTCPAddr("tcp", hostport)
	if err != nil { /* do something */ }
	conn, err := net.DialTCP("tcp", nil, tcpAddr) // maybe change nil to something
	if err != nil { /* do something */ }

	// create mapreduce request and write to server
	// change message type from string to different struct
	mr_request := []string{mrlib.MsgMAPREDUCE, file_name, start_line, end_line}
	mr_request_msg := strings.Join(mr_request, " ") 
	byte_mr_request_msg := []byte(mr_request_msg)
	err := conn.Write(byte_mr_request_msg)
	if err != nil { /* do something */ }


	// read answer from server
	byte_read_msg = [] // empty buffer, use some MAX SIZE?
	n, err := conn.Read(byte_read_msg)
	if err != nil { /* do something */ }

	// get message from byte_msg and print to command-line


	// if answer is good, print from file or something
	conn.Close()
	return
}