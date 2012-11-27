package main

import (
	"encoding/json"
	//"strconv"
	"mrlib"
	"net"
	"os"
	"os/exec"
	"bytes"
	"strconv"
	//"fmt"
)

const (
	verbosity = mrlib.Verbosity
)

func main() {

	// "./worker host:port"
	if len(os.Args) != 2 { return }

	hostport := os.Args[1]

	// Connect to server using TCP
	serverAddr, err := net.ResolveTCPAddr("tcp", hostport)
	if err != nil { /* do something */ }
	conn, err := net.DialTCP("tcp", nil, serverAddr) // maybe change nil to something
	if err != nil { /* do something */ }

	// identify as worker client
	identifyPacket := mrlib.IdentifyPacket{mrlib.MsgWORKERCLIENT}
	byteIdentifyPacket, err := json.Marshal(identifyPacket)
	if err != nil { /* do something */ }
	_ , err = conn.Write(byteIdentifyPacket)
	if err != nil { /* do something */ }

	for {
		// Read in Map or Reduce requests from server
		byteRequestMsg := make([]byte, mrlib.MaxMESSAGESIZE)
		n, err := conn.Read(byteRequestMsg[0:])
		if err != nil { /* do something */ }
		var request mrlib.ServerRequestPacket
		err = json.Unmarshal(byteRequestMsg[:n], &request)
		if err != nil { /* do something */ }

		answerPacket := mrlib.WorkerAnswerPacket{mrlib.MsgMAPANSWER, ""}

		switch (request.MsgType) {
		case mrlib.MsgMAPREQUEST:
			cmd := exec.Command(request.BinaryFile, "map",  request.FileName, strconv.Itoa(request.StartLine), strconv.Itoa(request.EndLine))
			var out bytes.Buffer 
			cmd.Stdout = &out 
			err := cmd.Start() 
			if err != nil {
				// log.Fatal(err)
			}
			err = cmd.Wait()
			// TODO : perform map job

			// send back results
			answerPacket.MsgType = mrlib.MsgMAPANSWER
			answerPacket.Answer = out.String() // TODO : change

		case mrlib.MsgREDUCEREQUEST:
			cmd := exec.Command(request.BinaryFile, "reduce",  request.FileName, strconv.Itoa(request.StartLine), strconv.Itoa(request.EndLine))
			var out bytes.Buffer 
			cmd.Stdout = &out 
			err := cmd.Start() 
			if err != nil {
				// log.Fatal(err)
			}
			err = cmd.Wait()
			// TODO : perform reduce job

			// send back results
			answerPacket.MsgType = mrlib.MsgREDUCEANSWER
			answerPacket.Answer = out.String() // TODO : change
		}

		// send answer back to server
		byteAnswer, err := json.Marshal(answerPacket)
		if err != nil { /* do something */ }
		n, err = conn.Write(byteAnswer)
		if err != nil { /* do something */ }


	}

}