package main

import (
	"os"
	"time"

	"github.com/lvhuat/http-service-example/testcmd/app"
	"github.com/lvhuat/http-service-example/testcmd/config"
	"github.com/lworkltd/kits/service/invoke"
	"github.com/sirupsen/logrus"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	cmd := kingpin.MustParse(app.App.Parse(os.Args[1:]))
	switch cmd {
	case config.Cmd.FullCommand():
		return
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: time.RFC3339,
	})

	config.Load()

	invoke.Init(&invoke.Option{
		DoLogger: !config.Config.Json,
	})

	var f func() error
	switch cmd {
	default:
	}

	f()

}
