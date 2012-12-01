package main

import (
	"container/list"
	"strconv"
	"mrlib"
	"bufio"
	"net"
	"os"
	"os/exec"
	"log"
	"fmt"
	"bytes"
	//"io"
)

type Request struct {
	conn *net.TCPConn
	mapJobs uint
	mapDone uint
	reduceJobs uint
	reduceDone uint
	binary string
	mapFile string
	output string
	input string
}

type Worker struct {
	conn *net.TCPConn
	job *mrlib.ServerRequestPacket
}

type mrServer struct {
	connListener *net.TCPListener
	workers map[uint]*Worker
	requests map[uint]*Request
	queue *list.List
	workerReady chan uint
	responseChan chan mrlib.WorkerAnswerPacket
	requestJoin chan uint
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
	server.eventHandler() // convert to go routine
}

func newServer(serverListener *net.TCPListener) *mrServer {
	server := &mrServer{}
	server.connListener = serverListener
	server.queue = list.New()
	server.workerReady = make(chan uint, 0)
	server.responseChan = make(chan mrlib.WorkerAnswerPacket, 0)
	server.requestJoin = make(chan uint, 0)
	server.requests = make(map[uint]*Request)
	server.workers = make(map[uint]*Worker)
	return server
}

func (server *mrServer) connectionHandler() {
	// only 2 iterations right now because 1 request/ 1 worker
	// later change to loop forever
	id := uint(1)
	for {
		conn, err := server.connListener.AcceptTCP()
		if err != nil { log.Fatal(err) }
		var identifyPacket mrlib.IdentifyPacket
		mrlib.Read(conn, &identifyPacket)
		log.Println(identifyPacket.MsgType)
		switch (identifyPacket.MsgType) {
		case mrlib.MsgREQUESTCLIENT:
			request := &Request{}
			request.conn = conn
			server.requests[id] = request
			go server.requestHandler(id, conn)
		case mrlib.MsgWORKERCLIENT:
			worker := &Worker{conn, nil}
			server.workers[id] = worker
			server.workerReady <- id
			go server.workerHandler(id, conn)
		}
		id++
	}
}

func (server *mrServer) eventHandler() {
	for {
		select {

		case id := <-server.requestJoin:
			server.addJobs(id)
			server.scheduleJobs()
			// Split job and add to queue
		case id := <-server.workerReady:
			server.workers[id].job = nil
			server.scheduleJobs()
		case answer := <-server.responseChan:
			size := answer.JobSize
			if answer.MsgType == mrlib.MsgMAPANSWER {
				f, err := os.OpenFile(server.requests[answer.RequestId].mapFile, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0666)
				f.WriteString(answer.Answer)
				log.Println(answer.Answer)
				f.Close()
				server.requests[answer.RequestId].mapDone += size
					// If all maps done, add requests to queue, sort file
				if server.requests[answer.RequestId].mapDone >= server.requests[answer.RequestId].mapJobs {
					cmd := exec.Command("sort", server.requests[answer.RequestId].mapFile)
					var out bytes.Buffer 
					cmd.Stdout = &out 
					err = cmd.Start()	
					if err != nil {
						log.Fatal("Sorting error: ", err)
					}
					err = cmd.Wait()
					f, _ := os.OpenFile(server.requests[answer.RequestId].mapFile, os.O_WRONLY, 0666)
					f.WriteString(out.String())
					f.Close()
					server.addJobs(answer.RequestId)
				}				
			} else {
				f, _ := os.OpenFile(server.requests[answer.RequestId].output, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0666)
				f.WriteString(answer.Answer)
				f.Close()
				server.requests[answer.RequestId].reduceDone += size
				if server.requests[answer.RequestId].reduceDone >= server.requests[answer.RequestId].reduceJobs {
					// If request done, remove map file, then output answer
					os.Remove(server.requests[answer.RequestId].mapFile)
					done := mrlib.MrAnswerPacket{mrlib.MsgSUCCESS}
					mrlib.Write(server.requests[answer.RequestId].conn, done)
				}
			} 
		}
	}
}

// reads in answers from worker clients
func (server *mrServer) workerHandler(id uint, conn *net.TCPConn) {

	var answer mrlib.WorkerAnswerPacket

	for {
		err := mrlib.Read(conn, &answer)
		if err != nil { log.Fatal(err) }
		server.responseChan <- answer
		server.workerReady <- id
	}

	return
}

// reads in requests from request clients
func (server *mrServer) requestHandler(id uint, conn *net.TCPConn) {

	fmt.Println("got request message")

	accept := mrlib.MrAnswerPacket{mrlib.MsgSUCCESS}
	mrlib.Write(conn, accept)

	var packet mrlib.MrRequestPacket
	mrlib.Read(conn, &packet)

	request := server.requests[id]
	request.binary = packet.BinaryFile
	request.input = packet.Directory
	request.mapFile = "tmp.txt" // generate unique file name here
	request.output = packet.AnswerFileName

	server.requests[id] = request

	server.requestJoin <- id
	return
}

func (server *mrServer) addJobs(id uint) {
	request := server.requests[id]
	// Have not added map jobs for this request yet
	if request.mapJobs == uint(0) {
		job := mrlib.ServerRequestPacket{}
		job.MsgType = mrlib.MsgMAPREQUEST
		job.FileName = request.input
		job.BinaryFile = request.binary

		file, err := os.Open(request.input) // TODO : change from file to dir
		if err != nil { }
		fileBuf := bufio.NewReader(file)
		
		startLine := 0
		endLine := 0

		// determine the length of the file
		err = nil
		for err == nil {
			_ , err = fileBuf.ReadString('\n')
			endLine++
		}

		chunk := mrlib.MrChunk{startLine, endLine} // do calcultions
		job.Ranges = []mrlib.MrChunk{chunk}
		job.RequestId = id
		job.JobSize = uint(1)

		server.queue.PushBack(job)

		request.mapJobs = uint(1)

	// Have not added reduce jobs for this request yet
	} else if request.reduceJobs == uint(0) {
		job := mrlib.ServerRequestPacket{}
		job.MsgType = mrlib.MsgREDUCEREQUEST
		job.FileName = request.mapFile
		job.BinaryFile = request.binary

		file, err := os.Open(request.mapFile) // TODO : change from file to dir
		if err != nil { }
		fileBuf := bufio.NewReader(file)
		
		startLine := 0
		endLine := 0

		// determine the length of the file
		err = nil
		for err == nil {
			_ , err = fileBuf.ReadString('\n')
			endLine++
		}

		chunk := mrlib.MrChunk{startLine, endLine-1} // do calcultions
		job.Ranges = []mrlib.MrChunk{chunk}
		job.RequestId = id
		job.JobSize = uint(1)

		server.queue.PushBack(job)

		request.reduceJobs = uint(1)
	}
}

func (server *mrServer) scheduleJobs() {
	for _, w := range(server.workers) {
		if server.queue.Len() <= 0 {
			return
		}
		if w.job == nil {
			job := server.queue.Remove(server.queue.Front()).(mrlib.ServerRequestPacket)
			mrlib.Write(w.conn, job)
		}
	}		
}