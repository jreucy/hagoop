package main

import (
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
	binaryName := os.Args[4]

	// connect to server with TCP
	serverAddr, err := net.ResolveTCPAddr("tcp", hostport)
	if err != nil { /* do something */ }
	conn, err := net.DialTCP("tcp", nil, serverAddr) // maybe change nil to something
	if err != nil { /* do something */ }

	// identify as request client
	identifyPacket := mrlib.IdentifyPacket{mrlib.MsgREQUESTCLIENT}
	mrlib.Write(conn, identifyPacket)

	// create mapreduce request and write to server
	mrRequest := mrlib.MrRequestPacket{fileDirectory, answerFileName, binaryName}
	mrlib.Write(conn, mrRequest)

	// read answer from server
	var answer mrlib.MrAnswerPacket
	answer = mrlib.Read(conn, answer).(mrlib.MrAnswerPacket)

	// get message from byte_msg and print to command-line
	switch (answer.MsgType) {
	case mrlib.MsgSUCCESS:
		// TODO : print from file or something
		return
	case mrlib.MsgFAIL:
		// fmt.Println("MapReduce failed")
		return
	}
	conn.Close()
	return
}