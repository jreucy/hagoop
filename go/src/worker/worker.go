package main

import (
	"strconv"
	"mrlib"
	"net"
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
	serverAddr, err := net.ResolveTCPAddr("tcp", hostport)
	if err != nil { /* do something */ }
	conn, err := net.DialTCP("tcp", nil, serverAddr) // maybe change nil to something
	if err != nil { /* do something */ }

	for {

		// Read in Map or Reduce requests from server
		byteRequestMsg := make([]byte, mrlib.MaxMESSAGESIZE)
		n, err := conn.Read(byteRequestMsg[0:])
		if err != nil { /* do something */ }
		var request mrlib.MrServerPacket
		err = json.Unmarshal(byteRequestMsg[:n], &request)
		if err != nil { /* do something */ }

		switch (request.MsgTYPE) {
		case mrlib.MsgMAPREQUEST:
			// perform map job

			// send back results
			mapAnswer := mrlib.MrWorkerPacket{mrlib.MsgMAPANSWER, ""}
			byteMapAnswer, err := json.Marshal(mapAnswer)
			if err != nil { /* do something */ }
			n, err = conn.Write(byteMapAnswer)
			if err != nil { /* do something */ }
		case mrlib.MsgREDUCEREQUEST:
			// perform reduce job

			// send back results
			reduceAnswer := mrlib.MrWorkerPacket{mrlib.MsgREDUCEANSWER, ""}
			byteReduceAnswer, err := json.Marshal(reduceAnswer)
			if err != nil { /* do something */ }
			n, err = conn.Write(byteMapAnswer)
			if err != nil { /* do something */ }
		}
	}

}