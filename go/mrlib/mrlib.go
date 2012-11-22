package mrlib

type packet struct {
	msgTYPE int
}

type mrfile struct {
	file string
	start_line int
	end_line int
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
)