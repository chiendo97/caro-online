package main

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/chiendo97/caro-online/internal/client"
	"github.com/chiendo97/caro-online/internal/game"
)

func run(ctx *cli.Context) error {
	var host = fmt.Sprintf("ws://%s:%d/find_hub", "localhost", 8080)
	var errC = make(chan error)

	type Worker struct {
		Id int
	}
	var queueSize = ctx.Int("worker")
	var workerQueue = make(chan Worker, queueSize)

	for i := 0; i < queueSize; i++ {
		var worker = Worker{i}
		workerQueue <- worker
	}

	logrus.Info("start")

	go func() {
		for worker := range workerQueue {
			go func(worker Worker) {
				defer func() {
					workerQueue <- worker
				}()

				logrus.Infof("worker %d start", worker.Id)

				c, _, err := websocket.DefaultDialer.Dial(host, nil)
				if err != nil {
					// logrus.Errorf("Dial error: %v", err)
					return
				}

				hub := client.NewHub(c, &game.RandomBot{})
				if err := hub.Run(); err != nil {
					logrus.Errorf("Hub run err: %v", err)
				}

				logrus.Infof("worker %d end", worker.Id)

			}(worker)
		}
	}()

	err := <-errC
	return err
}
