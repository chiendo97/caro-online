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
	logrus.SetFormatter(
		&logrus.TextFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return "", fmt.Sprintf("%-20s", fmt.Sprintf(" %s:%d ", filename, f.Line))
			},
		},
	)

	logrus.SetOutput(os.Stdout)

	logrus.SetLevel(logrus.InfoLevel)

	logrus.SetReportCaller(true)
}

func main() {

	app := &cli.App{
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

	err := app.Run(os.Args)
	if err != nil {
		logrus.Error(err)
	}
}
