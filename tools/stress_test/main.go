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
			return "", fmt.Sprintf(" %s:%d\t", filename, f.Line)
		},
	})

	logrus.SetOutput(os.Stdout)

	logrus.SetLevel(logrus.WarnLevel)

	logrus.SetReportCaller(true)
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
