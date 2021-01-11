package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/chiendo97/caro-online/internal/client"
	"github.com/chiendo97/caro-online/internal/game"
)

func run(ctx *cli.Context) error {
	var host = fmt.Sprintf("ws://%s:%d/find_hub", "localhost", 8080)
	var errC = make(chan error)

	go func() {
		for {
			time.Sleep(time.Second / 100)
			go func() {
				c, _, err := websocket.DefaultDialer.Dial(host, nil)
				if err != nil {
					logrus.Errorf("Dial error: %v", err)
					errC <- err
					return
				}

				hub := client.InitHub(c, &game.RandomBot{})
				if err := hub.Run(); err != nil {
					logrus.Errorf("Hub run err: %v", err)
				}
			}()
		}
	}()

	err := <-errC
	return err
}
