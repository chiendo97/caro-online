package socket

import (
	"errors"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type SocketI interface {
	SendMessage(msg Message)
	CloseMessage()
	RegisterHub(hub Hub)
}

type Socket struct {
	hub Hub

	conn *websocket.Conn

	msgC chan Message
}

func (s *Socket) SendMessage(msg Message) {
	s.msgC <- msg
}

func (s *Socket) CloseMessage() {
	close(s.msgC)
}

// InitAndRunSocket || xxx
func InitAndRunSocket(conn *websocket.Conn, hub Hub) *Socket {
	var socket = Socket{
		conn: conn,
		hub:  hub,
		msgC: make(chan Message),
	}

	go socket.read()
	go socket.write()

	return &socket
}

func (c *Socket) RegisterHub(hub Hub) {
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
		case msg, ok := <-c.msgC:

			if !ok {
				log.Info("socket: write closed")
				return
			}

			err := c.conn.WriteJSON(msg)
			if err != nil {
				log.Infof("socket: error write socket %v", err)
				return
			}
		}
	}
}

func (c *Socket) read() {
	defer func() {
		c.conn.Close()
		c.hub.UnRegister(c)
	}()

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Infof("socket: error read socket %v", err)
			} else {
				log.Info("socket: read closed")
			}
			return
		}

		go c.hub.HandleMsg(msg)
	}
}
