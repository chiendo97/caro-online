package main

import (
	"fmt"
	"sync"

	"github.com/chiendo97/caro-online/internal/client"
	"github.com/gorilla/websocket"
	"github.com/urfave/cli/v2"
)

func run(ctx *cli.Context) error {
	var host string

	host = fmt.Sprintf("ws://%s:%d/find_hub", "localhost", 8080)

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		// time.Sleep(time.Second / 2)
		go func() {
			defer wg.Done()
			c, _, err := websocket.DefaultDialer.Dial(host, nil)
			if err != nil {
				log.Errorf("Dial error: %v", err)
				return
			}

			hub := client.InitHub(c, &client.RandomBot{})

			if err := hub.Run(); err != nil {
				log.Errorf("Run: %v", err)
				return
			}
		}()
	}

	wg.Wait()

	return nil
}
