package main

import (
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/chiendo97/caro-online/internal/api"
	"github.com/chiendo97/caro-online/internal/server"
)

func run(c *cli.Context) error {
	core := server.InitCoreServer()

	service := api.InitService(core, c.Int("port"))

	errC := make(chan error)

	go func() {
		err := service.ListenAndServe(c.Int("port"))
		if err != nil {
			logrus.Errorf("Server run error: %v", err)
		}
		errC <- err
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	select {
	case err := <-errC:
		return err
	case <-interrupt:
		if err := service.Shutdown(); err != nil {
			logrus.Errorf("Shutdown server error: %v", err)
		}
		err := <-errC
		return err
	}
}
