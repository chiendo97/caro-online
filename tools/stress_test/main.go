package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/chiendo97/caro-online/cmd"
)

func main() {
	app := cmd.App{
		Action: run,
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "worker", Aliases: []string{"w"}, Value: 100},
		},
	}

	err := cmd.RunApp(app)
	if err != nil {
		logrus.Error(err)
	}
}
