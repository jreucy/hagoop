package mrlib

type MrRequestPacket struct {
	MsgTYPE int
	Directory string
	Answer_file_name string
}

type MrWorkerPacket struct {
	MsgTYPE int
	Answer string
}

type MrServerPacket struct {
	MsgType int
	File_name string
	Start_line int
	End_line int
}

type MrAnswerPacket struct {
	MsgType int
}

// maybe unnecessary
type MrFile struct {
	File string
	Start_line int
	End_line int
}

const (
  Verbosity = 1
  MsgJOIN = iota
  MsgMAPREDUCE
  MsgMAP_REQUEST
  MsgREDUCE_REQUEST
  MsgMAP_ANSWER
  MsgREDUCE_ANSWER
  MsgFAIL
  MsgSUCCESS
)