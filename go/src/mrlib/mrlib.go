package mrlib

type IdentifyPacket struct {
	MsgType int
}

type MrRequestPacket struct {
	Directory string
	AnswerFileName string
}

type WorkerAnswerPacket struct {
	MsgType int
	Answer string
}

type ServerRequestPacket struct {
	MsgType int
	FileName string
	StartLine int
	EndLine int
}

type MrAnswerPacket struct {
	MsgType int
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