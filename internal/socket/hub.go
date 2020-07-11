package socket

type Hub interface {
	// OnMessage is a function handling Msg
	OnMessage(msg Message)
	// UnRegister is a function handling when socket is closed
	UnRegister(s Socket)
}
