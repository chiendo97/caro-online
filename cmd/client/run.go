package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/chiendo97/caro-online/internal/client"
	"github.com/chiendo97/caro-online/internal/socket"
)

func run(ctx *cli.Context) error {
	var addr = ctx.String("addr")
	var port = ctx.Int("port")

	logrus.Printf("Client is connecting to %s:%d", addr, port)

	// === Take options
	var host string

	if ctx.String("join") != "" {
		host = fmt.Sprintf("ws://%s:%d/join_hub?hub=%s", addr, port, ctx.String("join"))
	} else if ctx.Bool("creat") {
		host = fmt.Sprintf("ws://%s:%d/create_hub", addr, port)
	} else if ctx.Bool("find") {
		host = fmt.Sprintf("ws://%s:%d/find_hub", addr, port)
	}

	if host == "" {
		return errors.New("xxx")
	}

	// === Init socket and hub
	c, _, err := websocket.DefaultDialer.Dial(host, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("Dial error: %v", err))
	}

	hub := client.InitAndRunHub()
	hub.Socket = socket.InitAndRunSocket(c, hub)

	// === take interrupt
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case <-interrupt:
			logrus.Fatalln("Exit client")
			return nil
		}
	}
}
