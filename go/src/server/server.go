package main

import (
	"container/list"
	"strconv"
	"mrlib"
	"bufio"
	"net"
	"os"
	//"fmt"
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
	binaryFileName string
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
	serverAddr, err := net.ResolveTCPAddr(mrlib.TCP, laddr)
	if err != nil { /* do something */ }
	serverListener, err := net.ListenTCP(mrlib.TCP, serverAddr) // maybe change nil to something
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
		if err != nil { /* do something */ }
		var identifyPacket mrlib.IdentifyPacket
		identifyPacket = mrlib.Read(conn, identifyPacket).(mrlib.IdentifyPacket)

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
				binaryFile := server.binaryFileName
				startLine := mrFile.StartLine
				endLine := mrFile.EndLine
				mapRequest := mrlib.ServerRequestPacket{ mrlib.MsgMAPREQUEST, mapFile, binaryFile, startLine, endLine }
				mrlib.Write(server.workerConn, mapRequest)
				// TODO : re-insert remaining file to front of list
			case <-server.reduceQueueNotEmpty:
				// send reduce request to next available worker
				reduceFile := ""
				binaryFile := server.binaryFileName
				startLine := 0
				endLine := 0
				reduceRequest := mrlib.ServerRequestPacket { mrlib.MsgREDUCEREQUEST, reduceFile, binaryFile, startLine, endLine }
				mrlib.Write(server.workerConn, reduceRequest)
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

	for {
		answer = mrlib.Read(server.workerConn, answer).(mrlib.WorkerAnswerPacket)

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
	var request mrlib.MrRequestPacket
	request = mrlib.Read(server.requestConn, request).(mrlib.MrRequestPacket)	

	server.binaryFileName = request.BinaryFile

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
	server.mapQueueNotEmpty <- true

	return
}