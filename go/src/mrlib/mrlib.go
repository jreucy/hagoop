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
	RequestId uint
	JobSize uint
	Answer string
}

type ServerRequestPacket struct {
	MsgType int
	FileName string
	BinaryFile string
	Ranges []MrChunk
	RequestId uint
	JobSize uint
}

type MrAnswerPacket struct {
	MsgType int
}

type MrJob struct {
	FileName string
	Ranges []MrChunk
}

type MrChunk struct {
	StartLine int
	EndLine int
}

func Write(conn *net.TCPConn, msg interface{}) error {
	byteMsg, err := json.Marshal(msg)
	if err != nil { return err }
	_ , err = conn.Write(byteMsg)
	if err != nil { return err }
	return nil
}

func Read(conn *net.TCPConn, varPointer interface{}) error {
	byteMsg := make([]byte, MaxMESSAGESIZE)
	n, err := conn.Read(byteMsg[0:])
	if err != nil { return err }
	err = json.Unmarshal(byteMsg[:n], varPointer)
	if err != nil { return err }
	return nil
}

func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}


const (
	Verbosity = 1
	TCP = "tcp"
	MAP = "map"
	REDUCE = "reduce"
	MaxMESSAGESIZE = 10000 // change later
	MinJOBSIZE = 100  // change
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