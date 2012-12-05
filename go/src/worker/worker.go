package main

import (
	"math/rand"
	"strconv"
	"os/exec"
	"mrlib"
	"bytes"
	"time"
	"log"
	"net"
	"os"
)

// specify chance of failure at every time interval
func fail(rate int) {
	for {
		time.Sleep(100 * time.Millisecond)
		n := rand.Intn(100)
		if rate > n {
			log.Println("Random Failure")
			os.Exit(0)
		}
	}
}

func main() {

	if len(os.Args) != 3 { return }

	hostport := os.Args[1]
	rate, _ := strconv.Atoi(os.Args[2])
	go fail(rate)

	// Connect to server using TCP as worker
	serverAddr, err := net.ResolveTCPAddr(mrlib.TCP, hostport)
	if err != nil { log.Fatal("Worker: ", err) }
	conn, err := net.DialTCP(mrlib.TCP, nil, serverAddr)
	if err != nil { log.Fatal("Worker: ", err) }
	identifyPacket := mrlib.IdentifyPacket{mrlib.MsgWORKERCLIENT}
	mrlib.Write(conn, identifyPacket)

	for {
		// Read in server requests and execute
		var request mrlib.ServerRequestPacket
		err = mrlib.Read(conn, &request)
		if err != nil { log.Fatal("Worker: ", err) }
		logJob(request)
		answerPacket := mrlib.WorkerAnswerPacket{}
		switch (request.MsgType) {
		case mrlib.MsgMAPREQUEST:
			ranges := request.Ranges
			var out bytes.Buffer 
			for i := 0; i < len(ranges); i++ {
				firstRange := ranges[i] 
				startLine := firstRange.StartLine
				endLine := firstRange.EndLine
				cmd := exec.Command(request.BinaryFile, mrlib.MAP, request.FileName, strconv.Itoa(startLine), strconv.Itoa(endLine), strconv.FormatInt(firstRange.Offset, 10))
				cmd.Stdout = &out 
				err := cmd.Start()	
				if err != nil { log.Fatal("Worker: ", err) }
				err = cmd.Wait()
			}
			answerPacket.MsgType = mrlib.MsgMAPANSWER
			answerPacket.Answer = out.String()
		case mrlib.MsgREDUCEREQUEST:
			ranges := request.Ranges		
			var out bytes.Buffer
			for i := 0; i < len(ranges); i++ {
				firstRange := ranges[i]
				startLine := firstRange.StartLine
				endLine := firstRange.EndLine
				cmd := exec.Command(request.BinaryFile, mrlib.REDUCE, request.FileName, strconv.Itoa(startLine), strconv.Itoa(endLine), strconv.FormatInt(firstRange.Offset, 10))
				cmd.Stdout = &out 
				err := cmd.Start() 
				if err != nil { log.Fatal("Worker: ", err) }
				err = cmd.Wait()
			}
			answerPacket.MsgType = mrlib.MsgREDUCEANSWER
			answerPacket.Answer = out.String()
		}
		// write answer back to the server
		answerPacket.JobSize = request.JobSize
		answerPacket.RequestId = request.RequestId
		mrlib.Write(conn, answerPacket)
	}
}

func logJob(request mrlib.ServerRequestPacket) {
	jobSize := 0
	ranges := request.Ranges
	for i := 0; i < len(ranges); i++ {
		r := ranges[i]
		startLine := r.StartLine
		endLine := r.EndLine
		jobSize += endLine - startLine
	}
	var msg string 
	if request.MsgType == mrlib.MsgREDUCEREQUEST {
		msg = "reduce"
	} else {
		msg = "map"
	}
	log.Println("Worker :", msg, "job size = ", jobSize)
}