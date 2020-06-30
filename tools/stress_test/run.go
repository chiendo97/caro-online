package main

import (
	"fmt"

	"github.com/chiendo97/caro-online/internal/client"
	"github.com/gorilla/websocket"
	"github.com/urfave/cli/v2"
)

func run(ctx *cli.Context) error {
	var host string

	host = fmt.Sprintf("ws://%s:%d/find_hub", "localhost", 8080)

	c, _, err := websocket.DefaultDialer.Dial(host, nil)
	if err != nil {
		return fmt.Errorf("Dial error: %v", err)
	}

	hub := client.InitHub(c, &client.RandomBot{})

	if err := hub.Run(); err != nil {
		return err
	}

	return nil
}
