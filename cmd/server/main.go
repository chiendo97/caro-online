package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return "", fmt.Sprintf(" %s:%d", filename, f.Line)
		},
	})

	log.SetOutput(os.Stdout)

	log.SetLevel(log.DebugLevel)

	log.SetReportCaller(true)
}

func main() {
	app := &cli.App{
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

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
}
