package main

import (
	"os"
	"os/signal"
	"sync"

	"github.com/chiendo97/caro-online/internal/api"
	"github.com/chiendo97/caro-online/internal/server"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func run(c *cli.Context) error {

	var wg sync.WaitGroup

	var core = server.InitCoreServer()

	var service = api.InitService(core, c.Int("port"))

	wg.Add(1)
	go func() {
		err := service.ListenAndServe(c.Int("port"))
		if err != nil {
			logrus.Errorf("Server run error: %v", err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := core.Run()
		if err != nil {
			logrus.Errorf("Core run error: %v", err)
		}
		wg.Done()
	}()

	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)
		<-interrupt
		core.Stop()
		if err := service.Shutdown(); err != nil {
			logrus.Errorf("Server shutdown error: %v", err)
		}
	}()

	wg.Wait()

	return nil
}
