package main

import (
	"encoding/json"
	"mrlib"
	"net"
	"os"
	//"fmt"
)

const (
	MAX_MESSAGE_SIZE = 1000 // change later
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
	conn, err := net.DialTCP("tcp", nil, serverAddr) // maybe change nil to something
	if err != nil { /* do something */ }

	// create mapreduce request and write to server
	// change message type from string to different struct

	mr_request := mrlib.MrRequestPacket{mrlib.MsgMAPREDUCE, file_directory, answer_file_name}
	byte_mr_request, err := json.Marshal(mr_request)
	if err != nil { /* do something */ }
	n, err := conn.Write(byte_mr_request)
	if err != nil { /* do something */ }


	// read answer from server
	byte_answer_msg := make([]byte, MAX_MESSAGE_SIZE)
	n, err = conn.Read(byte_answer_msg[0:])
	if err != nil { /* do something */ }


	// get message from byte_msg and print to command-line
	var answer mrlib.MrAnswerPacket
	err = json.Unmarshal(byte_answer_msg[:n], &answer)
	if err != nil { /* do something */ }
	switch (answer.MsgType) {
	// if answer is good, print from file or something
	case mrlib.MsgSUCCESS:
		break
	case mrlib.MsgFAIL:
		return
	default:
		return
	}
	conn.Close()
	return
}