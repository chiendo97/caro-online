package socket

type Hub interface {
	// OnMessage is a handle function when Msg comes
	OnMessage(msg Message)
	// OnLeave is a handle function when socket is closed
	OnLeave(s Socket)
}
