package socket

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Socket interface {
	Address() string
	SendMessage(msg interface{})
	Run() (error, error)
	Stop()
}

type socket struct {
	hub Hub

	conn *websocket.Conn

	msgC chan interface{}
	once sync.Once
}

func (s *socket) SendMessage(msg interface{}) {
	s.msgC <- msg
}

func (s *socket) Stop() {
	s.once.Do(func() {
		close(s.msgC)
	})
}

func NewSocket(conn *websocket.Conn, hub Hub) *socket {
	return &socket{
		conn: conn,
		hub:  hub,
		msgC: make(chan interface{}),
	}
}

func (c *socket) Run() (error, error) {
	defer c.conn.Close()

	logrus.Debugf("Socket %v start", c.Address())
	defer logrus.Debugf("Socket %v stop", c.Address())

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
	for {
		select {
		case msg, ok := <-c.msgC:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return nil
			}
			if err := c.conn.WriteJSON(msg); err != nil {
				if e, _ := err.(*websocket.CloseError); e != nil {
					logrus.Errorf("Write message err code: %v", e.Code)
				} else {
					logrus.Errorf("Write message err: %v", err)
				}
				return err
			}
		}
	}
}

func (c *socket) read() error {
	defer c.hub.OnLeave(c)

	var msg Message
	for {
		if err := c.conn.ReadJSON(&msg); err != nil {
			if websocket.IsCloseError(
				err,
				websocket.CloseNoStatusReceived,
				websocket.CloseNormalClosure,
			) {
				return nil
			}
			if e, _ := err.(*websocket.CloseError); e != nil {
				logrus.Errorf("Read message err code: %v", e.Code)
			} else {
				logrus.Errorf("Read message err: %v", err)
			}
			return err
		}

		c.hub.OnMessage(msg)
	}
}

func (c *socket) Address() string {
	return c.conn.RemoteAddr().String()
}
