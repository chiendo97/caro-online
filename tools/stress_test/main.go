package main

import (
	"github.com/sirupsen/logrus"

	"github.com/chiendo97/caro-online/cmd"
)

func main() {
	app := cmd.App{
		Action: run,
	}

	err := cmd.RunApp(app)
	if err != nil {
		logrus.Error(err)
	}
}
