package main

import (
	"mrlib"
	"net"
	"os"
	"os/exec"
	"bytes"
	"strconv"
	"log"
	"fmt"
)

const (
	verbosity = mrlib.Verbosity
)

func main() {

	// "./worker host:port"
	if len(os.Args) != 2 { return }

	hostport := os.Args[1]

	// Connect to server using TCP
	serverAddr, err := net.ResolveTCPAddr(mrlib.TCP, hostport)
	if err != nil { log.Fatal("Worker: ", err) }
	conn, err := net.DialTCP(mrlib.TCP, nil, serverAddr) // maybe change nil to something
	if err != nil { log.Fatal("Worker: ", err) }

	// identify as worker client
	identifyPacket := mrlib.IdentifyPacket{mrlib.MsgWORKERCLIENT}
	mrlib.Write(conn, identifyPacket)

	for {
		// Read in Map or Reduce requests from server
		var request mrlib.ServerRequestPacket
		err = mrlib.Read(conn, &request)
		if err != nil { 
			log.Fatal("Worker: ", err)
		}
		answerPacket := mrlib.WorkerAnswerPacket{}
		switch (request.MsgType) {
		case mrlib.MsgMAPREQUEST:
			ranges := request.Ranges
			firstRange := ranges[0] // assumes single chunk, put in for loop later
			startLine := firstRange.StartLine
			endLine := firstRange.EndLine
			cmd := exec.Command(request.BinaryFile, mrlib.MAP,  request.FileName, strconv.Itoa(startLine), strconv.Itoa(endLine))
			var out bytes.Buffer 
			cmd.Stdout = &out 
			err := cmd.Start()	
			if err != nil { log.Fatal("Worker: ", err) }
			err = cmd.Wait()

			// send back results
			answerPacket.JobSize = request.JobSize
			answerPacket.RequestId = request.RequestId
			answerPacket.MsgType = mrlib.MsgMAPANSWER
			answerPacket.Answer = out.String() // TODO : change
			fmt.Println("Received map message")

		case mrlib.MsgREDUCEREQUEST:
			fmt.Println("Received reduce message")
			ranges := request.Ranges			
			firstRange := ranges[0] // assumes single chunk, put in for loop later
			startLine := firstRange.StartLine
			endLine := firstRange.EndLine
			cmd := exec.Command(request.BinaryFile, mrlib.REDUCE,  request.FileName, strconv.Itoa(startLine), strconv.Itoa(endLine))
			var out bytes.Buffer 
			cmd.Stdout = &out 
			err := cmd.Start() 
			if err != nil {
				log.Fatal(err)
			}
			err = cmd.Wait()

			// send back results
			answerPacket.JobSize = request.JobSize
			answerPacket.RequestId = request.RequestId
			answerPacket.MsgType = mrlib.MsgREDUCEANSWER
			answerPacket.Answer = out.String() // TODO : change

		}
		mrlib.Write(conn, answerPacket)
	}

}