package socket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Hub interface {
	ReceiveMsg(msg Message)
	Unregister(s *Socket)
}

type Socket struct {
	// Hub for behide execution
	Hub Hub

	// The websocket connection.
	Conn *websocket.Conn

	// receive Message from hub and send thorough Conn
	Message chan Message
}

func InitSocket(conn *websocket.Conn, hub Hub) *Socket {
	var socket = Socket{
		Conn:    conn,
		Hub:     hub,
		Message: make(chan Message),
	}

	go socket.Read()
	go socket.Write()

	return &socket
}

func (c *Socket) Write() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Message:

			if !ok {
				log.Println("socket: write closed")
				// c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.Conn.WriteJSON(msg)
			if err != nil {
				log.Printf("socket: error write socket %v", err)
				return
			}
		}
	}
}

func (c *Socket) Read() {
	defer func() {
		c.Conn.Close()
		c.Hub.Unregister(c)
	}()

	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("socket: error read socket %v", err)
			} else {
				log.Println("socket: read closed")
			}
			return
		}

		go c.Hub.ReceiveMsg(msg)
	}
}
