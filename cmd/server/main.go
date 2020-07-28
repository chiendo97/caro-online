package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/chiendo97/caro-online/cmd"
)

func main() {
	app := cmd.App{
		Name:  "caro-online",
		Usage: "run caro-online server",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Usage:       "server `PORT`",
				DefaultText: "8080",
				Value:       8080,
			},
		},
		Action: run,
	}

	err := cmd.RunApp(app)
	if err != nil {
		logrus.Error(err)
	}
}
