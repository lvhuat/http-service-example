package main

import (
	"cryptobroker/tradenode/testcmd/testtrade/app"
	"cryptobroker/tradenode/testcmd/testtrade/config"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/lworkltd/kits/service/invoke"

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
