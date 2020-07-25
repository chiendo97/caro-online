package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/chiendo97/caro-online/internal/client"
	"github.com/chiendo97/caro-online/internal/game"
)

func run(ctx *cli.Context) error {

	var (
		addr = ctx.String("addr")
		port = ctx.Int("port")
		host string
	)

	logrus.Printf("Client is connecting to %s:%d", addr, port)

	// === Take options
	switch ctx.String("option") {
	case "find":
		host = fmt.Sprintf("ws://%s:%d/find_hub", addr, port)
	case "join":
		host = fmt.Sprintf("ws://%s:%d/join_hub?hub=%s", addr, port, ctx.String("gameID"))
	case "create":
		host = fmt.Sprintf("ws://%s:%d/create_hub", addr, port)
	default:
		return fmt.Errorf("Invalid option: %s", ctx.String("option"))
	}

	if host == "" {
		return fmt.Errorf("")
	}

	// === Init socket and hub
	c, _, err := websocket.DefaultDialer.Dial(host, nil)
	if err != nil {
		return fmt.Errorf("Dial error: %v", err)
	}

	hub := client.InitHub(c, &game.RandomBot{})

	errC := make(chan error)

	go func() {
		errC <- hub.Run()
	}()

	// === take interrupt
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case <-interrupt:
			hub.Stop()
		case err := <-errC:
			logrus.Info("Exit client")
			return err
		}
	}
}
