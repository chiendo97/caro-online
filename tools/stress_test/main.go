package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return "", fmt.Sprintf(" %s:%d\t", filename, f.Line)
		},
	})

	log.SetOutput(os.Stdout)

	log.SetLevel(logrus.ErrorLevel)

	log.SetReportCaller(true)
}

func main() {
	fmt.Println("stress test")

	app := &cli.App{
		Action: run,
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Error(err)
	}
}
