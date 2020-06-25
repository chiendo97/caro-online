package main

import (
	"github.com/chiendo97/caro-online/internal/api"
	"github.com/chiendo97/caro-online/internal/server"
	"github.com/urfave/cli/v2"
)

func run(c *cli.Context) error {
	var core = server.InitCoreServer()
	go func() {
		core.Run()
	}()

	var service = api.InitService(core)

	err := service.ListenAndServe(c.Int("port"))
	if err != nil {
		return err
	}

	return nil
}
