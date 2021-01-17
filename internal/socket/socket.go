package socket

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Socket interface {
	GetSocketIPAddress() string
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
	exporterCounter.WithLabelValues("write").Inc()
	s.msgC <- msg
}

func (s *socket) Stop() {
	s.once.Do(func() {
		close(s.msgC)
	})
}

func NewSocket(conn *websocket.Conn, hub Hub) *socket {
	var s = socket{
		conn: conn,
		hub:  hub,
		msgC: make(chan interface{}),
	}

	return &s
}

func (c *socket) Run() (error, error) {
	defer c.conn.Close()

	logrus.Debugf("Socket %v start", c.GetSocketIPAddress())
	defer logrus.Debugf("Socket %v stop", c.GetSocketIPAddress())

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
	defer c.conn.Close()

	for msg := range c.msgC {
		start := time.Now()
		err := c.conn.WriteJSON(msg)
		exporterLatency.WithLabelValues("write").Observe(float64(time.Since(start).Milliseconds()))
		if err != nil {
			if e, _ := err.(*websocket.CloseError); e != nil {
				logrus.Errorf("Write message err code: %v", e.Code)
			} else {
				logrus.Errorf("Write message err: %v", err)
			}
			return err
		}
	}

	return nil
}

func (c *socket) read() error {
	defer c.hub.OnLeave(c)

	var msg Message
	for {
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				return nil
			}
			if e, _ := err.(*websocket.CloseError); e != nil {
				logrus.Errorf("Read message err code: %v", e.Code)
			} else {
				logrus.Errorf("Read message err: %v", err)
			}
			return err
		}

		go c.hub.OnMessage(msg)
	}
}

func (c *socket) GetSocketIPAddress() string {
	return c.conn.RemoteAddr().String()
}
