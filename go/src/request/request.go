package main

import (
	"mrlib"
	"net"
	"log"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) != 5 { return }

	hostport := os.Args[1]
	fileDirectory := os.Args[2]
	answerFileName := os.Args[3]
	binaryName := os.Args[4]

	// connect to server with TCP as request client
	serverAddr, err := net.ResolveTCPAddr(mrlib.TCP, hostport)
	if err != nil { log.Fatal("Request: ", err) }
	conn, err := net.DialTCP(mrlib.TCP, nil, serverAddr)
	if err != nil { log.Fatal("Request: ", err) }
	identifyPacket := mrlib.IdentifyPacket{mrlib.MsgREQUESTCLIENT}
	mrlib.Write(conn, identifyPacket)

	var acceptPacket mrlib.MrAnswerPacket
	mrlib.Read(conn, &acceptPacket)
	if acceptPacket.MsgType != mrlib.MsgSUCCESS { log.Fatal("Request not accepted") }

	// create mapreduce request and write to server
	mrRequest := mrlib.MrRequestPacket{fileDirectory, answerFileName, binaryName}
	mrlib.Write(conn, mrRequest)

	// read answer from server
	var answer mrlib.MrAnswerPacket
	mrlib.Read(conn, &answer)

	switch (answer.MsgType) {
	case mrlib.MsgSUCCESS:
		fmt.Println("Success! Result saved in", answerFileName)
		return
	case mrlib.MsgFAIL:
		fmt.Println("MapReduce failed")
		return
	}
	conn.Close()
	return
}