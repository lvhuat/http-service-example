package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/lvhuat/http-service-example/testcmd/app"
	"github.com/lworkltd/kits/service/restful/code"
	"github.com/lworkltd/kits/utils/jsonize"
	"github.com/sirupsen/logrus"
)

type XFeignKey struct {
	TradePlatform  string `json:"tradePlatform"`
	TradeSite      string `json:"tradeSite"`
	TradeAccessId  string `json:"tradeAccessKey"`
	TradeSecretKey string `json:"tradeSecretKey"`
	TradeAccountId string `json:"tradeAccountId"`
	UserId         string `json:"userId"`
}

func (me *XFeignKey) Json() string {
	return jsonize.V(me, false)
}

type ServerRouter struct {
	Remote   string     `json:"remote"`
	Scheme   string     `json:"scheme"`
	Json     bool       `json:"json"`
	FeignKey *XFeignKey `json:"feignKey"`
}

var Config ServerRouter

var (
	Cmd        = app.App.Command("set", "建议手动直接修改config.json文件")
	configFile = app.App.Flag("config_file", "Load config file").Short('f').Default("config.json").String()
	scheme     = Cmd.Flag("scheme", "Config the http scheme for request using").Required().Enum("http", "https")
	remote     = Cmd.Flag("remote", "Config the http host&port for request using").Required().String()
)

func Excute() {
	config := &ServerRouter{
		Scheme: *scheme,
		Remote: *remote,
	}
	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(*configFile, b, 0777); err != nil {
		panic(err)
	}

	fmt.Println("Set config done.")
}

func Load() {
	b, err := ioutil.ReadFile(*configFile)
	if err != nil {
		fmt.Println("Please check config.json exist")
		os.Exit(-1)
	}

	if err := json.Unmarshal(b, &Config); err != nil {
		fmt.Println("bad config format")
		os.Exit(-1)
	}

	if Config.Json {
		logrus.SetLevel(logrus.FatalLevel)
	}
}

var since = time.Now()

func CostMs() int64 {
	return int64(time.Now().Sub(since) / time.Millisecond)
}

type Profile struct {
	Method  string
	Path    string
	Input   interface{}
	Output  interface{}
	Code    string
	Message string
	DelayMs int64
}

func DealJson(method, url string, req, rsp interface{}, err code.Error) {
	if Config.Json {
		profile := Profile{
			Path:    url,
			Method:  "POST",
			DelayMs: CostMs(),
		}
		if err == nil {
			profile.Input = req
			profile.Output = rsp
		} else {
			profile.Code = err.Mcode()
			profile.Message = err.Message()
		}

		fmt.Println(jsonize.V(profile, true))
	}
}
