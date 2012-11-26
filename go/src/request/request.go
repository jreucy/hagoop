package main

import (
	"encoding/json"
	"mrlib"
	"net"
	"os"
	//"fmt"
)

const (
	MaxMESSAGESIZE = 10000 // change later
)

func main() {

	// "./request host:port [file_directory] [file_name]"
	if len(os.Args) != 4 { return }

	hostport := os.Args[1]
	fileDirectory := os.Args[2]
	answerFileName := os.Args[3]

	// connect to server with TCP
	serverAddr, err := net.ResolveTCPAddr("tcp", hostport)
	if err != nil { /* do something */ }
	conn, err := net.DialTCP("tcp", nil, serverAddr) // maybe change nil to something
	if err != nil { /* do something */ }

	// create mapreduce request and write to server
	// change message type from string to different struct

	mrRequest := mrlib.MrRequestPacket{mrlib.MsgMAPREDUCE, fileDirectory, answerFileName}
	byteMrRequest, err := json.Marshal(mrRequest)
	if err != nil { /* do something */ }
	n, err := conn.Write(byteMrRequest)
	if err != nil { /* do something */ }


	// read answer from server
	byteAnswerMsg := make([]byte, MaxMESSAGESIZE)
	n, err = conn.Read(byteAnswerMsg[0:])
	if err != nil { /* do something */ }


	// get message from byte_msg and print to command-line
	var answer mrlib.MrAnswerPacket
	err = json.Unmarshal(byteAnswerMsg[:n], &answer)
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