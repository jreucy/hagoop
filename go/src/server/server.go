package main

import (
	"encoding/json"
	"strconv"
	"mrlib"
	"net"
	//"fmt"
	"os"
)

const (
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

	// connect to server with TCP
	serverAddr, err := net.ResolveTCPAddr("tcp", laddr)
	if err != nil { /* do something */ }
	serverListener, err := net.ListenTCP("tcp", serverAddr) // maybe change nil to something
	if err != nil { /* do something */ }

	// assuming a single request
	// assuming request comes first
	// later put in own goroutine
	requestConn, err := serverListener.AcceptTCP()
	if err != nil { /* do something */ }

	// assuming worker comes second
	workerConn, err := serverListener.AcceptTCP()
	if err != nil { /* do something */ }

	// read in request packet
	byteRequestMsg := make([]byte, mrlib.MaxMESSAGESIZE)
	n, err := requestConn.Read(byteRequestMsg[0:])
	if err != nil { /* do something */ }
	var request mrlib.MrRequestPacket
	err = json.Unmarshal(byteRequestMsg[:n], &request)
	if err != nil { /* do something */ }
	if request.MsgTYPE != mrlib.MsgMAPREDUCE { return /* or do something else */ }
	//directory := request.Directory
	//answerFileName := request.AnswerFileName

	// parse directory and save file name, starting/ending line numbers
	mapFile := ""
	startLine := 0
	endLine := 0

	// send map request
	mapRequest := mrlib.MrServerPacket{ mrlib.MsgMAPREQUEST, mapFile, startLine, endLine }
	byteMapRequest, err := json.Marshal(mapRequest)
	if err != nil { /* do something */ }
	n, err = workerConn.Write(byteMapRequest)
	if err != nil { /* do something */ }

	// read map answer
	byteAnswerMsg := make([]byte, mrlib.MaxMESSAGESIZE)
	n, err = workerConn.Read(byteAnswerMsg[0:])
	if err != nil { /* do something */ }
	var answer mrlib.MrWorkerPacket
	err = json.Unmarshal(byteRequestMsg[:n], &answer)
	if err != nil { /* do something */ }
	if answer.MsgTYPE != mrlib.MsgMAPANSWER { /* do something */ }
	//mapAnswer := answer.Answer

	// save string into file
	reduceFile := ""
	startLine = 0
	endLine = 0

	// send read request
	reduceRequest := mrlib.MrServerPacket { mrlib.MsgREDUCEREQUEST, reduceFile, startLine, endLine }
	byteReduceRequest, err := json.Marshal(reduceRequest)
	if err != nil { /* do something */ }
	n, err = workerConn.Write(byteReduceRequest)
	if err != nil { /* do something */ }

	// read reduce answer
	byteAnswerMsg = make([]byte, mrlib.MaxMESSAGESIZE)
	n, err = workerConn.Read(byteAnswerMsg[0:])
	if err != nil { /* do something */ }
	err = json.Unmarshal(byteRequestMsg[:n], &answer)
	if err != nil { /* do something */ }
	if answer.MsgTYPE != mrlib.MsgREDUCEANSWER { /* do something */ }
	//reduceAnswer := answer.Answer

	// write reduce answer to specified file

	// send mapreduce answer to request





	// for when we have multiple workers/clients
	//for {
		// Read in incoming requests from both request and worker clients



		// schedule and assign pending jobs to available workers

	//}
}