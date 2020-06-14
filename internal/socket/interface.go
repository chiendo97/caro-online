package socket

type Hub interface {
	// HandleMsg is a function handling Msg
	HandleMsg(msg Message)
	// UnRegister is a function handling when socket is closed
	UnRegister(s *Socket)
}
