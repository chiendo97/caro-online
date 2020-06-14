package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return "", fmt.Sprintf(" %s:%d", filename, f.Line)
		},
	})

	logrus.SetOutput(os.Stdout)

	logrus.SetLevel(logrus.DebugLevel)

	logrus.SetReportCaller(true)
}

func main() {

	app := &cli.App{
		Name:  "caro-online",
		Usage: "run caro-online client",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "find",
				Usage:   "Find Game",
				Aliases: []string{"f"},
				Value:   true,
			},
			&cli.BoolFlag{
				Name:    "create",
				Usage:   "Create Game",
				Aliases: []string{"c"},
				Value:   false,
			},
			&cli.StringFlag{
				Name:        "join",
				Usage:       "Join game ID",
				Aliases:     []string{"j"},
				DefaultText: "",
				Value:       "",
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

	err := app.Run(os.Args)
	if err != nil {
		logrus.Error(err)
	}
}
