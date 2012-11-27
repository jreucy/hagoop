package main

import (
	"container/list"
	"encoding/json"
	"strconv"
	"mrlib"
	"bufio"
	"net"
	//"fmt"
	"os"
	//"io"
)

const (
	//MIN_JOB_SIZE = 1000  // change
	//MAX_JOB_SIZE = 10000 // change
)

type mrServer struct {
	connListener *net.TCPListener
	requestConn *net.TCPConn      // single request client
	workerConn *net.TCPConn       // single worker client
	mapList *list.List
	mapQueueNotEmpty chan bool    // placeholder, change to something better later
	reduceQueueNotEmpty chan bool // look above
	finishedAllMaps chan bool
	finishedAllReduces chan bool
	saveMapToFile chan string
	saveReduceToFile chan string
}

func main() {

	// "./server port"
	if len(os.Args) != 2 { return }

	port := os.Args[1]

	// make sure port is an integer
	_ , err := strconv.Atoi(port)
	if err != nil { /* do something */ }

	// connect to server with TCP
	laddr := ":" + port
	serverAddr, err := net.ResolveTCPAddr("tcp", laddr)
	if err != nil { /* do something */ }
	serverListener, err := net.ListenTCP("tcp", serverAddr) // maybe change nil to something
	if err != nil { /* do something */ }

	server := newServer(serverListener)
	go server.connectionHandler()
	go server.eventHandler()
}

func newServer(serverListener *net.TCPListener) *mrServer {
	server := &mrServer{}
	server.connListener = serverListener
	server.requestConn = nil
	server.workerConn = nil
	server.mapList = list.New()
	server.mapQueueNotEmpty = make(chan bool, 0)
	server.reduceQueueNotEmpty = make(chan bool, 0)
	server.finishedAllMaps = make(chan bool, 0)
	server.finishedAllReduces = make(chan bool, 0)
	server.saveMapToFile = make(chan string, 0)
	server.saveReduceToFile = make(chan string, 0)
	return server
}

func (server *mrServer) connectionHandler() {
	// only 2 iterations right now because 1 request/ 1 worker
	// later change to loop forever
	for i := 0; i < 2; i++ {
		conn, err := server.connListener.AcceptTCP()
		byteMsg := make([]byte, mrlib.MaxMESSAGESIZE)
		n, err := conn.Read(byteMsg[0:])
		if err != nil { /* do something */ }
		var identifyPacket mrlib.IdentifyPacket
		err = json.Unmarshal(byteMsg[:n], &identifyPacket)
		if err != nil { /* do something */ }

		switch (identifyPacket.MsgType) {
		case mrlib.MsgREQUESTCLIENT:
			server.requestConn = conn
			go server.requestHandler()
		case mrlib.MsgWORKERCLIENT:
			server.workerConn = conn
			go server.workerHandler()
		default:
			// do something, break, etc.
		}
	}
}

func (server *mrServer) eventHandler() {
	for {
		select {
			case <-server.mapQueueNotEmpty:
				// send map request to next available worker
				mrFile := server.mapList.Remove(server.mapList.Front()).(mrlib.MrFile)
				mapFile := mrFile.FileName
				startLine := mrFile.StartLine
				endLine := mrFile.EndLine
				mapRequest := mrlib.ServerRequestPacket{ mrlib.MsgMAPREQUEST, mapFile, startLine, endLine }
				byteMapRequest, err := json.Marshal(mapRequest)
				if err != nil { /* do something */ }
				_ , err = server.workerConn.Write(byteMapRequest)
				if err != nil { /* do something */ }
				// TODO : re-insert remaining file to front of list
			case <-server.reduceQueueNotEmpty:
				// send reduce request to next available worker
				reduceFile := ""
				startLine := 0
				endLine := 0
				reduceRequest := mrlib.ServerRequestPacket { mrlib.MsgREDUCEREQUEST, reduceFile, startLine, endLine }
				byteReduceRequest, err := json.Marshal(reduceRequest)
				if err != nil { /* do something */ }
				_ , err = server.workerConn.Write(byteReduceRequest)
				if err != nil { /* do something */ }
			case <-server.saveMapToFile:
				// save map answer to file
			case <-server.saveReduceToFile:
				// write reduce answer to specified file
			case <-server.finishedAllMaps:
				// put request jobs in request queue 
			case <-server.finishedAllReduces:
				// send mapreduce answer to request client
		}
	}
}

// reads in answers from worker clients
func (server *mrServer) workerHandler() {

	var answer mrlib.WorkerAnswerPacket
	byteAnswerMsg := make([]byte, mrlib.MaxMESSAGESIZE)

	for {
		n, err := server.workerConn.Read(byteAnswerMsg[0:])
		if err != nil { /* do something */ }
		err = json.Unmarshal(byteAnswerMsg[:n], &answer)
		if err != nil { /* do something */ }

		switch (answer.MsgType) {
		case mrlib.MsgMAPANSWER:
			server.saveMapToFile <- answer.Answer
		case mrlib.MsgREDUCEANSWER:
			server.saveReduceToFile <- answer.Answer
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

	// parse directory and save file name, starting/ending line numbers
	// currently "directory" represents a single file
	file, err := os.Open(request.Directory) // TODO : change from directory to file
	if err != nil { /* do something */ }
	fileBuf := bufio.NewReader(file)
	
	startLine := 0
	endLine := 0

	// determine the length of the file
	err = nil
	for err == nil {
		_ , err = fileBuf.ReadString('\n')
		endLine++
	}
	mrFile := mrlib.MrFile{request.Directory, startLine, endLine}

	// place map job into buffer
	server.mapList.PushBack(mrFile)

	return
}