package socket

import (
	"fmt"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Socket interface {
	GetSocketIPAddress() string
	SendMessage(msg Message)
	CloseMessage()
	Run() (error, error)
}

type socket struct {
	hub Hub

	conn *websocket.Conn

	msgC chan Message

	done chan struct{}
}

func (s *socket) SendMessage(msg Message) {
	s.msgC <- msg
}

func (s *socket) CloseMessage() {
	close(s.msgC)
}

func InitSocket(conn *websocket.Conn, hub Hub) *socket {
	var s = socket{
		conn: conn,
		hub:  hub,
		msgC: make(chan Message),
		done: make(chan struct{}),
	}

	return &s
}

func (c *socket) Run() (error, error) {

	defer func() {
		c.conn.Close()
	}()

	if c.hub == nil {
		return fmt.Errorf("Hub is missing"), fmt.Errorf("Hub is missing")
	}

	errC := make(chan error)

	go func() {
		errC <- c.read()
	}()

	go func() {
		errC <- c.write()
	}()

	return <-errC, <-errC
}

// GetSocketIPAddress returns ip address of socket
func (c *socket) GetSocketIPAddress() string {
	return c.conn.RemoteAddr().String()
}

func (c *socket) write() error {

	for {
		select {
		case <-c.done:
			return nil
		case msg, ok := <-c.msgC:

			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					return err
				}
				return nil
			}

			err := c.conn.WriteJSON(msg)
			if err != nil {
				log.Infof("socket: error write socket %v", err)
				return err
			}
		}
	}
}

func (c *socket) read() error {
	defer func() {
		close(c.done)
		go c.hub.UnRegister(c)
	}()

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)

		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				return nil
			}
			e, _ := err.(*websocket.CloseError)
			log.Infof("socket: error read socket %v %v", err, e)
			return err
		}

		go c.hub.HandleMsg(msg)
	}
}
