package main

import (
	"container/list"
	"strconv"
	"mrlib"
	"bufio"
	"net"
	"os"
	"log"
	"fmt"
	"time"
	//"io"
)

const (
	WorkerFREE = iota
	WorkerBUSY
)

type Worker struct {
	status int
	conn *net.TCPConn
}

type mrServer struct {
	connListener *net.TCPListener
	requestConn *net.TCPConn      // single request client  
	worker Worker 				  // single worker client
	mapList *list.List
	reduceList *list.List
	binaryFileName string
	mapAnswerFile string 		  // Single map output file
	answerFileName string 		  // Single output file
	mapQueueNotEmpty chan bool    // placeholder, change to something better later
	reduceQueueNotEmpty chan bool // look above
	finishedAllMaps chan bool
	finishedAllReduces chan bool
	saveMapToFile chan string
	saveReduceToFile chan string
	changeWorkerStatus chan bool
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
	if err != nil { log.Fatal(err) }
	serverListener, err := net.ListenTCP(mrlib.TCP, serverAddr)
	if err != nil { log.Fatal(err) }

	server := newServer(serverListener)
	go server.connectionHandler()
	go server.MapQueueHandler()
	server.eventHandler() // convert to go routine
}

func newServer(serverListener *net.TCPListener) *mrServer {
	server := &mrServer{}
	server.connListener = serverListener
	server.requestConn = nil
	server.worker = Worker{}
	server.mapAnswerFile = "tmp.txt"
	server.mapList = list.New()
	server.reduceList = list.New()
	server.mapQueueNotEmpty = make(chan bool, 0)
	server.reduceQueueNotEmpty = make(chan bool, 0)
	server.finishedAllMaps = make(chan bool, 0)
	server.finishedAllReduces = make(chan bool, 0)
	server.saveMapToFile = make(chan string, 0)
	server.saveReduceToFile = make(chan string, 0)
	server.changeWorkerStatus = make(chan bool, 0)
	return server
}

func (server *mrServer) connectionHandler() {
	// only 2 iterations right now because 1 request/ 1 worker
	// later change to loop forever
	for i := 0; i < 2; i++ {
		conn, err := server.connListener.AcceptTCP()
		if err != nil { log.Fatal(err) }
		var identifyPacket mrlib.IdentifyPacket
		mrlib.Read(conn, &identifyPacket)
		log.Println(identifyPacket.MsgType)
		switch (identifyPacket.MsgType) {
		case mrlib.MsgREQUESTCLIENT:
			server.requestConn = conn
			go server.requestHandler()
		case mrlib.MsgWORKERCLIENT:
			server.worker.status = WorkerFREE
			server.worker.conn = conn
			go server.workerHandler()
		// default:
			// do something, break, etc.
		}
	}
}

func (server *mrServer) eventHandler() {
	for {
		select {
			case <-server.changeWorkerStatus:
				server.worker.status = WorkerFREE
			case <-server.mapQueueNotEmpty:
				// send map request to next available worker
				mrJob := server.mapList.Remove(server.mapList.Front()).(mrlib.MrJob)
				mapFile := mrJob.FileName
				binaryFile := server.binaryFileName
				ranges := mrJob.Ranges
				mapRequest := mrlib.ServerRequestPacket{ mrlib.MsgMAPREQUEST, mapFile, binaryFile, ranges }
				mrlib.Write(server.worker.conn, mapRequest)
				server.worker.status = WorkerBUSY
			case <-server.reduceQueueNotEmpty:
				// send reduce request to next available worker
				mrJob := server.reduceList.Remove(server.reduceList.Front()).(mrlib.MrJob)
				reduceFile := server.mapAnswerFile
				binaryFile := server.binaryFileName
				ranges := mrJob.Ranges
				reduceRequest := mrlib.ServerRequestPacket { mrlib.MsgREDUCEREQUEST, reduceFile, binaryFile, ranges }
				mrlib.Write(server.worker.conn, reduceRequest)
			case answer := <-server.saveMapToFile:
				file, err := os.Create(server.mapAnswerFile)
				if err != nil { log.Fatal(err) }
				file.WriteString(answer)
				file.Close()
				server.saveMapToFile <- "done"
			case answer := <-server.saveReduceToFile:
				file, err := os.Create(server.answerFileName)
				if err != nil { log.Fatal(err) }
				file.WriteString(answer)
				file.Close()
				// write reduce answer to specified file
			case <-server.finishedAllMaps:
				// put request jobs in request queue 
			case <-server.finishedAllReduces:
				os.Remove(server.mapAnswerFile)
				// Remove mapAnswer file now that reduces are done
				done := mrlib.MrAnswerPacket{mrlib.MsgSUCCESS}
				mrlib.Write(server.requestConn, done)
				return // send mapreduce answer to request client
		}
	}
}

// reads in answers from worker clients
func (server *mrServer) workerHandler() {

	var answer mrlib.WorkerAnswerPacket

	for {
		mrlib.Read(server.worker.conn, &answer)
		switch (answer.MsgType) {
		case mrlib.MsgMAPANSWER:
			server.saveMapToFile <- answer.Answer
			<-server.saveMapToFile

			file, err := os.Open(server.mapAnswerFile)
			if err != nil { log.Fatal(err) }
			fileBuf := bufio.NewReader(file)
			
			startLine := 0
			endLine := 0

			// determine the length of the file
			err = nil
			for err == nil {
				_ , err = fileBuf.ReadString('\n')
				endLine++
			}
			chunk := mrlib.MrChunk{startLine, endLine - 1} // Overcounts by 1
			ranges := []mrlib.MrChunk{chunk}
			mrFile := mrlib.MrJob{server.mapAnswerFile, ranges}	

			// place reduce job into buffer
			server.reduceList.PushBack(mrFile)
			file.Close()
			server.reduceQueueNotEmpty <- true

		case mrlib.MsgREDUCEANSWER:
			server.saveReduceToFile <- answer.Answer

			server.finishedAllReduces <- true
		}
	}

	return
}

// reads in requests from request clients
func (server *mrServer) requestHandler() {

	fmt.Println("got request message")

	accept := mrlib.MrAnswerPacket{mrlib.MsgSUCCESS}
	mrlib.Write(server.requestConn, accept)
	// put in for loop later for multiple request clients

	// read in request packet
	var request mrlib.MrRequestPacket
	mrlib.Read(server.requestConn, &request)	

	server.binaryFileName = request.BinaryFile
	server.answerFileName = request.AnswerFileName

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
	/* for when we actually split the jobs
	for i := startLine; i < endLine; i += mrlib.MinJOBSIZE {
		jobStart := i
		jobEnd := mrlib.Min(endLine, i + mrlib.MinJOBSIZE)
		chunk := mrlib.MrChunk{jobStart, jobEnd}
		ranges := []mrlib.MrChunk{chunk}
		mrJob := mrlib.MrJob{request.Directory, ranges}
		server.mapList.PushBack(mrJob)
	}
	*/
	chunk := mrlib.MrChunk{startLine, endLine}
	ranges := []mrlib.MrChunk{chunk}
	mrJob := mrlib.MrJob{request.Directory, ranges}
	server.mapList.PushBack(mrJob)
	return
}

/* definitely fix */
func (server *mrServer) MapQueueHandler() {
	for {
		if server.mapList.Len() > 0 {
			// assuming single worker
			if server.worker.status == WorkerFREE {
				server.mapQueueNotEmpty <- true
			}
		} 			
		time.Sleep(500 * time.Millisecond)
	}
}