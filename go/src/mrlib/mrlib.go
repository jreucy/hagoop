package mrlib

import (
	"encoding/json"
	"net"
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
	if err != nil { /* do something */ }
	_ , err = conn.Write(byteMsg)
	if err != nil { /* do something */ }
}

func Read(conn *net.TCPConn, varPointer interface{}) interface{} {
	byteMsg := make([]byte, MaxMESSAGESIZE)
	readStruct := varPointer
	n, err := conn.Read(byteMsg[0:])
	if err != nil { /* do something */ }
	err = json.Unmarshal(byteMsg[:n], &readStruct)
	if err != nil { /* do something */ }
	return readStruct
}


const (
	Verbosity = 1
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