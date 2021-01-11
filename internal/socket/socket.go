package socket

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Socket interface {
	// GetSocketIPAddress returns ip address of socket
	GetSocketIPAddress() string

	SendMessage(msg Message)

	Run() (error, error)

	CloseMessage()
}

type socket struct {
	hub Hub

	conn *websocket.Conn

	isClosed bool
	msgC     chan Message

	once sync.Once
	mux  sync.Mutex
}

func (s *socket) SendMessage(msg Message) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.isClosed {
		return
	}

	s.msgC <- msg
}

func (s *socket) CloseMessage() {
	s.once.Do(func() {
		s.mux.Lock()
		defer s.mux.Unlock()

		s.isClosed = true
		close(s.msgC)
	})
}

func InitSocket(conn *websocket.Conn, hub Hub) *socket {
	var s = socket{
		conn:     conn,
		hub:      hub,
		msgC:     make(chan Message),
		isClosed: false,
	}

	return &s
}

func (c *socket) Run() (error, error) {

	logrus.Debugf("Socket %v start", c.GetSocketIPAddress())
	defer logrus.Debugf("Socket %v stop", c.GetSocketIPAddress())

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

	defer func() {
		c.conn.Close()
	}()

	for msg := range c.msgC {
		err := c.conn.WriteJSON(msg)
		if err != nil {
			e, _ := err.(*websocket.CloseError)
			if e != nil {
				logrus.Warnf("Write message err code: %v", e.Code)
			}
			logrus.Warnf("Write message err: %v", err)
			return err
		}
	}

	// if msgC is closed by `CloseMessage`
	c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	return nil
}

func (c *socket) read() error {
	defer func() {
		go c.hub.OnLeave(c)
	}()

	var msg Message
	for {
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				return nil
			}
			e, _ := err.(*websocket.CloseError)
			logrus.Errorf("socket: error read socket %v %v", err, e)
			return err
		}

		go c.hub.OnMessage(msg)
	}
}

func (c *socket) GetSocketIPAddress() string {
	return c.conn.RemoteAddr().String()
}
