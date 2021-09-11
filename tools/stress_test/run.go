package main

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/chiendo97/caro-online/internal/client"
	"github.com/chiendo97/caro-online/internal/game"
)

type Worker struct {
	Id int
}

func run(ctx *cli.Context) error {
	host := fmt.Sprintf("ws://%s:%d/find_hub", "localhost", 8080)

	queueSize := ctx.Int("worker")
	workerQueue := make(chan Worker, queueSize)

	for i := 0; i < queueSize; i++ {
		worker := Worker{i}
		workerQueue <- worker
	}

	logrus.Info("start")

	for i := 0; i < queueSize; i++ {
		go func() {
			for {
				worker := <-workerQueue

				logrus.Infof("worker %d start", worker.Id)

				c, _, err := websocket.DefaultDialer.Dial(host, nil)
				if err != nil {
					return
				}

				hub := client.NewHub(c, &game.RandomBot{})
				if err := hub.Run(); err != nil {
					logrus.Errorf("Hub run err: %v", err)
				}

				logrus.Infof("worker %d end", worker.Id)

				workerQueue <- worker
			}
		}()
	}

	select {}
}
