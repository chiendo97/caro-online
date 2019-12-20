package socket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Socket struct {
	// Hub for behide execution
	Hub interface {
		ReceiveMsg(msg Message)
		Unregister(s *Socket)
	}

	// The websocket connection.
	Conn *websocket.Conn

	// receive Message from hub and send thorough Conn
	Message chan Message
}

func (c *Socket) Write() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Message:
			if !ok {
				log.Println("Message chan close")
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.Conn.WriteJSON(msg)

			if err != nil {
				log.Println("Cant send msg:", err, msg)
				return
			}
		}
	}
}

func (c *Socket) Read() {
	defer func() {
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()

	for {
		var msg Message

		err := c.Conn.ReadJSON(&msg)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Cant read msg: %v", err)
			} else {
				log.Println("I read close message")
			}
			return
		}

		c.Hub.ReceiveMsg(msg)
	}
}
