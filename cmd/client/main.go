package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/chiendo97/caro-online/cmd"
)

func main() {
	app := cmd.App{
		Name:  "caro-online",
		Usage: "run caro-online client",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "option",
				Usage:       "find | join | create",
				Aliases:     []string{"o"},
				Value:       "find",
				DefaultText: "find",
			},
			&cli.StringFlag{
				Name:    "gameID",
				Usage:   "gameID for joinning game",
				Aliases: []string{"g"},
			},
			&cli.StringFlag{
				Name:        "addr",
				Usage:       "Host location",
				DefaultText: "localhost",
				Value:       "localhost",
			},
			&cli.IntFlag{
				Name:        "port",
				Usage:       "Port",
				Aliases:     []string{"p"},
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
