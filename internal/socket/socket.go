package socket

import (
	"errors"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Socket interface {
	GetSocketIPAddress() string
	SendMessage(msg Message)
	CloseMessage()
}

type socket struct {
	hub Hub

	conn *websocket.Conn

	msgC chan Message
}

func (s *socket) SendMessage(msg Message) {
	s.msgC <- msg
}

func (s *socket) CloseMessage() {
	close(s.msgC)
}

// InitAndRunSocket || xxx
func InitAndRunSocket(conn *websocket.Conn, hub Hub) *socket {
	var s = socket{
		conn: conn,
		hub:  hub,
		msgC: make(chan Message),
	}

	go s.read()
	go s.write()

	return &s
}

func (c *socket) Run() error {

	if c.hub == nil {
		return errors.New("Hub is missing")
	}

	go c.read()
	go c.write()

	return nil
}

// GetSocketIPAddress returns ip address of socket
func (c *socket) GetSocketIPAddress() string {
	return c.conn.RemoteAddr().String()
}

func (c *socket) write() {
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

func (c *socket) read() {
	defer func() {
		c.hub.UnRegister(c)
		c.conn.Close()
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
