package mrlib

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