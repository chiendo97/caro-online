package cmd

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type App struct {
	Name  string
	Usage string

	Flags  []cli.Flag
	Action cli.ActionFunc
}

var commonFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "log",
		Aliases: []string{"l"},
		Usage:   "debug | info | error | warm",
		Value:   "debug",
	},
}

func RunApp(a App) error {
	flags := append(commonFlags, a.Flags...)

	app := &cli.App{
		Name:   a.Name,
		Usage:  a.Usage,
		Flags:  flags,
		Action: run(a.Action),
	}

	return app.Run(os.Args)
}

func run(runFunc cli.ActionFunc) cli.ActionFunc {
	return func(c *cli.Context) error {
		if err := initLog(c.String("log")); err != nil {
			return err
		}

		if err := runFunc(c); err != nil {
			return err
		}

		return nil
	}
}

func initLog(levelStr string) error {
	logrus.SetFormatter(
		&logrus.TextFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return "", fmt.Sprintf("%-20s", fmt.Sprintf(" %s:%d ", filename, f.Line))
			},
		},
	)

	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)

	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)

	return nil
}
