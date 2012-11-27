package main

import (
	"mrlib"
	"net"
	"os"
	"os/exec"
	"bytes"
	"strconv"
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
	mrlib.Write(conn, identifyPacket)

	for {
		// Read in Map or Reduce requests from server
		var request mrlib.ServerRequestPacket
		request = mrlib.Read(conn, request).(mrlib.ServerRequestPacket)

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

		mrlib.Write(conn, answerPacket)
	}

}