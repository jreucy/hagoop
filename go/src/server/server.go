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

type mrServer struct {
	requestConn *net.TCPConn      // single request client
	workerConn *net.TCPConn       // single worker client
	mapQueueNotEmpty chan bool    // placeholder, change to something better later
	reduceQueueNotEmpty chan bool // look above
	finishedAllMaps chan bool
	finishedAllReduces chan bool
}

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

	// create new server object
	server := &mrServer{}

	// put into go routine to accept incoming connections
	// only 2 iterations now because 1 request/ 1 worker
	for i := 0; i < 2; i++ {
		conn, err := serverListener.AcceptTCP()
		byteMsg := make([]byte, mrlib.MaxMESSAGESIZE)
		n, err := conn.Read(byteMsg[0:])
		if err != nil { /* do something */ }
		var requestPacket mrlib.MrRequestPacket
		var workerPacket mrlib.MrWorkerPacket
		errRequest := json.Unmarshal(byteMsg[:n], &requestPacket)
		errWorker := json.Unmarshal(byteMsg[:n], &workerPacket)
		if errRequest == nil {
			server.requestConn = conn
			go server.requestHandler()
		} else if errWorker == nil {
			server.workerConn = conn
			go server.workerHandler()
		} else {
			// do something
		}
	}

	go server.eventHandler()
}

func (server *mrServer) eventHandler() {
	for {
		select {
			case <-server.mapQueueNotEmpty:
				// send map request
				mapFile := ""
				startLine := 0
				endLine := 0
				mapRequest := mrlib.MrServerPacket{ mrlib.MsgMAPREQUEST, mapFile, startLine, endLine }
				byteMapRequest, err := json.Marshal(mapRequest)
				if err != nil { /* do something */ }
				_ , err = server.workerConn.Write(byteMapRequest)
				if err != nil { /* do something */ }
			case <-server.reduceQueueNotEmpty:
				// send reduce request 
				reduceFile := ""
				startLine := 0
				endLine := 0
				reduceRequest := mrlib.MrServerPacket { mrlib.MsgREDUCEREQUEST, reduceFile, startLine, endLine }
				byteReduceRequest, err := json.Marshal(reduceRequest)
				if err != nil { /* do something */ }
				_ , err = server.workerConn.Write(byteReduceRequest)
				if err != nil { /* do something */ }
			case <-server.finishedAllMaps:
				// put request jobs in request queue 
			case <-server.finishedAllReduces:
				// send mapreduce answer to request client
		}
	}
}

// reads in answers from worker clients
func (server *mrServer) workerHandler() {

	var answer mrlib.MrWorkerPacket
	byteAnswerMsg := make([]byte, mrlib.MaxMESSAGESIZE)

	for {
		n, err := server.workerConn.Read(byteAnswerMsg[0:])
		if err != nil { /* do something */ }
		err = json.Unmarshal(byteAnswerMsg[:n], &answer)
		if err != nil { /* do something */ }

		switch (answer.MsgTYPE) {
		case mrlib.MsgMAPANSWER:
			// mapAnswer := answer.Answer
			// save string into file
			break
		case mrlib.MsgREDUCEANSWER:
			// reduceAnswer := answer.Answer
			// write reduce answer to specified file
			break
		}
	}

	return
}

// reads in requests from request clients
func (server *mrServer) requestHandler() {

	// put in for loop later for multiple request clients

	// read in request packet
	byteRequestMsg := make([]byte, mrlib.MaxMESSAGESIZE)
	n, err := server.requestConn.Read(byteRequestMsg[0:])
	if err != nil { /* do something */ }
	var request mrlib.MrRequestPacket
	err = json.Unmarshal(byteRequestMsg[:n], &request)
	if err != nil { /* do something */ }
	if request.MsgTYPE != mrlib.MsgMAPREDUCE { return /* or do something else */ }
	//directory := request.Directory
	//answerFileName := request.AnswerFileName

	// parse directory and save file name, starting/ending line numbers

	// place map jobs into buffer

	return
}