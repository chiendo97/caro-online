package socket

import (
	"errors"
	"log"

	"github.com/gorilla/websocket"
)

type hubG interface {
	ReceiveMsg(msg Message)
	Unregister(s *Socket)
}

type Socket struct {
	hub hubG

	conn *websocket.Conn

	// receive Message from hub and send thorough Conn
	Message chan Message
}

// InitAndRunSocket || xxx
func InitAndRunSocket(conn *websocket.Conn, hub hubG) *Socket {
	var socket = Socket{
		conn:    conn,
		hub:     hub,
		Message: make(chan Message),
	}

	go socket.read()
	go socket.write()

	return &socket
}

func (c *Socket) RegisterHub(hub hubG) {
	c.hub = hub
}

func (c *Socket) Run() error {

	if c.hub == nil {
		return errors.New("Hub is missing")
	}

	go c.read()
	go c.write()

	return nil
}

// GetSocketIPAddress returns ip address of socket
func (c *Socket) GetSocketIPAddress() string {
	return c.conn.RemoteAddr().String()
}

func (c *Socket) write() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Message:

			if !ok {
				log.Println("socket: write closed")
				return
			}

			err := c.conn.WriteJSON(msg)
			if err != nil {
				log.Printf("socket: error write socket %v", err)
				return
			}
		}
	}
}

func (c *Socket) read() {
	defer func() {
		c.conn.Close()
		c.hub.Unregister(c)
	}()

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("socket: error read socket %v", err)
			} else {
				log.Println("socket: read closed")
			}
			return
		}

		go c.hub.ReceiveMsg(msg)
	}
}