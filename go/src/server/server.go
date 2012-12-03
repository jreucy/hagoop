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
	"bytes"
	"math"
	"time"
	"io/ioutil"
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
	id uint
	conn *net.TCPConn
	job *mrlib.ServerRequestPacket
	joinTime time.Time
	writeChan chan *mrlib.ServerRequestPacket
}

type mrServer struct {
	connListener *net.TCPListener
	workers map[uint]*Worker
	requests map[uint]*Request
	queue *list.List
	workerReady chan *Worker
	workerDied chan uint
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
	server.workerReady = make(chan *Worker, 0)
	server.workerDied = make(chan uint, 0)
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
		switch (identifyPacket.MsgType) {
		case mrlib.MsgREQUESTCLIENT:
			request := &Request{}
			request.conn = conn
			server.requests[id] = request
			go server.requestHandler(id, conn)
		case mrlib.MsgWORKERCLIENT:
			worker := &Worker{id, conn, nil, time.Now(), make(chan *mrlib.ServerRequestPacket, 0)}
			server.workerReady <- worker
			go server.workerHandler(worker)
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
		case worker := <-server.workerReady:
			worker.job = nil
			server.workers[worker.id] = worker
			server.scheduleJobs()
		case id := <-server.workerDied:
			worker := server.workers[id]
			if worker.job != nil {
				server.splitFailedJob(worker.job)
			}
			delete(server.workers, id)
			server.scheduleJobs()
		case answer := <-server.responseChan:
			size := answer.JobSize
			if answer.MsgType == mrlib.MsgMAPANSWER {
				f, err := os.OpenFile(server.requests[answer.RequestId].mapFile, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0666)
				f.WriteString(answer.Answer)
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
func (server *mrServer) workerHandler(worker *Worker) {

	log.Println("Worker joined with id:", worker.id)

	var answer mrlib.WorkerAnswerPacket

	for {
		job := <- worker.writeChan
		mrlib.Write(worker.conn, job)
		worker.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		err := mrlib.Read(worker.conn, &answer)
		if err != nil {
			log.Println("Read error: Worker", worker.id, "died")
			worker.conn.SetLinger(0)
			worker.conn.Close()
			server.workerDied <- worker.id
			return
		}
		server.responseChan <- answer
		server.workerReady <- worker
	}

	return
}

// reads in requests from request clients
func (server *mrServer) requestHandler(id uint, conn *net.TCPConn) {

	log.Println("Request received with id:", id)

	accept := mrlib.MrAnswerPacket{mrlib.MsgSUCCESS}
	err := mrlib.Write(conn, accept)
	if err != nil {
		log.Println("Write Error: Request", id, "died")
		return
	}

	var packet mrlib.MrRequestPacket
	err = mrlib.Read(conn, &packet)
	if err != nil {
		log.Println("Read Error: Request", id, "died")
		return
	}

	request := server.requests[id]
	request.binary = packet.BinaryFile
	request.input = packet.Directory
	request.mapFile = ".mapFile" + strconv.FormatUint(uint64(id), 10)
	request.output = packet.AnswerFileName

	server.requests[id] = request

	server.requestJoin <- id
	return
}

func (server *mrServer) splitFailedJob(job *mrlib.ServerRequestPacket) {
	ranges := job.Ranges

	for i := 0; i < len(ranges); i++ {
		newServerRequest := &mrlib.ServerRequestPacket{}
		newServerRequest.MsgType = job.MsgType
		newServerRequest.FileName = job.FileName
		newServerRequest.BinaryFile = job.BinaryFile
		newServerRequest.Ranges = []mrlib.MrChunk{ranges[i]}
		newServerRequest.RequestId = job.RequestId
		newServerRequest.JobSize = 1
		server.queue.PushBack(newServerRequest)
	}
}

func (server *mrServer) addDir(id uint, dirName string) {
	request := server.requests[id]
	files, _ := ioutil.ReadDir(dirName)
	for f := 0; f < len(files); f++ {
		if files[f].IsDir() {
			server.addDir(id, dirName + "/" + files[f].Name())
		}
		lines := countLines(dirName + "/" + files[f].Name())
		chunks := splitMapJob(lines)
		for i := 0; i < len(chunks); i++ {
			job := &mrlib.ServerRequestPacket{}
			job.MsgType = mrlib.MsgMAPREQUEST
			job.FileName = dirName + "/" + files[f].Name()
			job.BinaryFile = request.binary
			job.Ranges = []mrlib.MrChunk{chunks[i]}
			job.RequestId = id
			job.JobSize = uint(1)

			server.queue.PushBack(job)
		}

		request.mapJobs += uint(len(chunks))
	}
}

func (server *mrServer) addJobs(id uint) {
	request := server.requests[id]
	// Have not added map jobs for this request yet
	if request.mapJobs == uint(0) {
		server.addDir(id, request.input)

	// Have not added reduce jobs for this request yet
	} else if request.reduceJobs == uint(0) {
		
		lines := countLines(request.mapFile)
		chunks := splitReduceJob(lines, request.mapFile)
		for i := 0; i < len(chunks); i++ {
			job := &mrlib.ServerRequestPacket{}
			job.MsgType = mrlib.MsgREDUCEREQUEST
			job.FileName = request.mapFile
			job.BinaryFile = request.binary
			job.Ranges = []mrlib.MrChunk{chunks[i]}
			job.RequestId = id
			job.JobSize = uint(1)

			server.queue.PushBack(job)
		}

		request.reduceJobs = uint(len(chunks))
	}
}

func (server *mrServer) scheduleJobs() {
	for id, w := range(server.workers) {
		if server.queue.Len() <= 0 {
			return
		}
		if w.job == nil {
			job := server.combineJobs(w)
			w.writeChan <- job
			w.job = job
			server.workers[id] = w
		}
	}		
}

func (server *mrServer) combineJobs(worker *Worker) *mrlib.ServerRequestPacket {
	workerLife := int(time.Now().Sub(worker.joinTime) / (1 * time.Second))
	maxJobSize := mrlib.Min(workerLife, mrlib.MaxJOBNUM)
	firstJob := server.queue.Remove(server.queue.Front()).(*mrlib.ServerRequestPacket)
	jobsToRemove := make([]*list.Element, mrlib.MaxJOBNUM)
	i := 0
	for job := server.queue.Front(); job != nil && i < maxJobSize; job = job.Next() {
		j := job.Value.(*mrlib.ServerRequestPacket)
		sameType := j.MsgType == firstJob.MsgType
		sameFile := j.FileName == firstJob.FileName
		sameRequest := j.RequestId == firstJob.RequestId
		if sameType && sameFile && sameRequest {
			firstRange := firstJob.Ranges
			jRange := j.Ranges
			firstRange = append(firstRange, jRange...)
			firstJob.Ranges = firstRange
			firstJob.JobSize++
			jobsToRemove[i] = job
			i++
		}
	}

	for x := 0; x < i; x++ {
		server.queue.Remove(jobsToRemove[x])
	}

	// implement combining here
	return firstJob
}

func splitMapJob(numLines int) []mrlib.MrChunk {
	// split into chunks based on min job size
	numJobs := int(math.Ceil(float64(numLines)/float64(mrlib.MinJOBSIZE)))
	chunks := make([]mrlib.MrChunk, numJobs)
	start := 0
	end := mrlib.MinJOBSIZE - 1
	for i := 0; i < numJobs; i++ {
		if end > numLines {
			end = numLines
		}
		chunks[i] = mrlib.MrChunk{start, end}
		start = end
		end += mrlib.MinJOBSIZE
	}
	return chunks
}

func splitReduceJob(numLines int, mapFile string) []mrlib.MrChunk {

	maxJobs := int(math.Ceil(float64(numLines)/float64(mrlib.MinJOBSIZE)))
	chunks := make([]mrlib.MrChunk, maxJobs)

	file, err := os.Open(mapFile)
	if err != nil { }
	fileBuf := bufio.NewReader(file)
	
	// loop through file, finding number of unique keys
	// save start and end lines of each unique key
	i := 1
	numUniqueKeys := 0
	keyStart := 0
	keyEnd := 0
	firstKeyArr, err := fileBuf.ReadString('\n')
	firstKey := mrlib.GetKey(firstKeyArr)
	err = nil
	for err == nil {
		keyArr, err := fileBuf.ReadString('\n')
		if err != nil { 
			chunks[numUniqueKeys] = mrlib.MrChunk{keyStart, i}	
			numUniqueKeys++
			break 
		}
		key := mrlib.GetKey(keyArr)

		if key != firstKey {

			firstKey = key
			if i - keyStart >= mrlib.MinJOBSIZE {
				keyEnd = i
				chunks[numUniqueKeys] = mrlib.MrChunk{keyStart, keyEnd}
				keyStart = keyEnd
				numUniqueKeys++
			}
		}
		i++
	}

	// Same keys have to be on same worker
	return chunks[0:numUniqueKeys]
}

func countLines(fileName string) int {
	file, err := os.Open(fileName)
	if err != nil { }
	fileBuf := bufio.NewReader(file)
	
	lines := 0

	// determine the length of the file
	err = nil
	for err == nil {
		lines++
		_ , err = fileBuf.ReadString('\n')
	}
	return lines
}