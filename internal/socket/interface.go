package socket

type hubG interface {
	// ReceiveMsg is a function handling Msg
	ReceiveMsg(msg Message)
	// Unregister is a function handling when socket is closed
	Unregister(s *Socket)
}
