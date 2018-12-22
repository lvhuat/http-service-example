package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lvhuat/http-service-example/conf"
	"github.com/lvhuat/http-service-example/dto"
	"github.com/lvhuat/http-service-example/httpserver"

	"runtime"

	"github.com/Sirupsen/logrus"
)

func perr(action string, err error) {
	logrus.WithFields(logrus.Fields{
		"error": err,
	}).Fatalln(action, "Start server failed!!!")

	os.Exit(-1)
}

func main() {
	runtime.GOMAXPROCS(2)

	logrus.SetLevel(logrus.DebugLevel)
	if err := conf.Parse(); err != nil {
		perr("conf.Parse", err)
	}
	conf.Dump()

	if err := dto.Init(conf.GetMysql().Url); err != nil {
		perr("dto.Init", err)
	}

	if err := httpserver.Run(conf.GetService()); err != nil {
		perr("httpserver.Run", err)
	}
}
