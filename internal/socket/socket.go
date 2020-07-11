package socket

import (
	"fmt"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Socket interface {
	GetSocketIPAddress() string
	SendMessage(msg Message)
	Run() (error, error)
	CloseMessage()
}

type socket struct {
	hub Hub

	conn *websocket.Conn

	msgC chan Message
}

func (s *socket) SendMessage(msg Message) {
	log.Warnf("Sending socket(%v) msg(%v)", s.GetSocketIPAddress(), msg)
	s.msgC <- msg
	log.Warnf("Sending done socket(%v) msg(%v)", s.GetSocketIPAddress(), msg)
}

func (s *socket) CloseMessage() {
	close(s.msgC)
}

func InitSocket(conn *websocket.Conn, hub Hub) *socket {
	var s = socket{
		conn: conn,
		hub:  hub,
		msgC: make(chan Message),
	}

	return &s
}

func (c *socket) Run() (error, error) {

	log.Debugf("Socket %v start", c.GetSocketIPAddress())
	defer log.Debugf("Socket %v stop", c.GetSocketIPAddress())

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

func (c *socket) write() error {

	log.Debugf("Socket write %v start", c.GetSocketIPAddress())
	defer log.Debugf("Socket write %v stop", c.GetSocketIPAddress())

	for {
		select {
		case msg, ok := <-c.msgC:

			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return nil
			}

			err := c.conn.WriteJSON(msg)
			if err != nil {
				e, _ := err.(*websocket.CloseError)
				if e != nil {
					log.Warnf("Write message err code: %v", e.Code)
				}
				log.Warnf("Write message err: %v", err)
			}
		}
	}
}

func (c *socket) read() error {
	defer func() {
		go c.hub.UnRegister(c)
	}()

	log.Debugf("Socket read %v start", c.GetSocketIPAddress())
	defer log.Debugf("Socket read %v stop", c.GetSocketIPAddress())

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

// GetSocketIPAddress returns ip address of socket
func (c *socket) GetSocketIPAddress() string {
	return c.conn.RemoteAddr().String()
}
