package mrlib

import (
	"encoding/json"
	"net"
	"log"
)


type IdentifyPacket struct {
	MsgType int
}

type MrRequestPacket struct {
	Directory string
	AnswerFileName string
	BinaryFile string
}

type WorkerAnswerPacket struct {
	MsgType int
	Answer string
}

type ServerRequestPacket struct {
	MsgType int
	FileName string
	BinaryFile string
	StartLine int
	EndLine int
}

type MrAnswerPacket struct {
	MsgType int
}

type MrFile struct {
	FileName string
	StartLine int
	EndLine int
}

func Write(conn *net.TCPConn, msg interface{}) {
	byteMsg, err := json.Marshal(msg)
	if err != nil { log.Fatal("Write error: ", err) }
	_ , err = conn.Write(byteMsg)
	if err != nil { log.Fatal("Write error: ", err)  }
}

func Read(conn *net.TCPConn, varPointer interface{}) {
	byteMsg := make([]byte, MaxMESSAGESIZE)
	n, err := conn.Read(byteMsg[0:])
	if err != nil { log.Fatal("Read error: ", err) }
	err = json.Unmarshal(byteMsg[:n], varPointer)
	if err != nil { log.Fatal("Read error: ", err) }
}


const (
	Verbosity = 1
	TCP = "tcp"
	MAP = "map"
	REDUCE = "reduce"
	MaxMESSAGESIZE = 10000 // change later
	MinJOBSIZE = 1000  // change
	MaxJOBSIZE = 10000 // change
	MsgREQUESTCLIENT = iota
	MsgWORKERCLIENT
	MsgMAPREQUEST
	MsgREDUCEREQUEST
	MsgMAPANSWER
	MsgREDUCEANSWER
	MsgFAIL
	MsgSUCCESS
)