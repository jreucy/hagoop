package mrlib

type MrRequestPacket struct {
	MsgTYPE int
	Directory string
	AnswerFileName string
}

type MrWorkerPacket struct {
	MsgTYPE int
	Answer string
}

type MrServerPacket struct {
	MsgType int
	FileName string
	StartLine int
	EndLine int
}

type MrAnswerPacket struct {
	MsgType int
}

// maybe unnecessary
type MrFile struct {
	File string
	StartLine int
	EndLine int
}

const (
	Verbosity = 1
	MaxMESSAGESIZE = 10000 // change later
	MinJOBSIZE = 1000  // change
	MaxJOBSIZE = 10000 // change
	MsgJOIN = iota
	MsgMAPREDUCE
	MsgMAPREQUEST
	MsgREDUCEREQUEST
	MsgMAPANSWER
	MsgREDUCEANSWER
	MsgFAIL
	MsgSUCCESS
)