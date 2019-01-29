package conf

import (
	"fmt"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/lvhuat/http-service-example/common"
	"github.com/lworkltd/kits/service/discovery"
	"github.com/lworkltd/kits/service/invoke"
	"github.com/lworkltd/kits/service/profile"
	"github.com/lworkltd/kits/utils/jsonize"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type Profile struct {
	Base        profile.Base
	Service     profile.Service
	Invoker     profile.Invoker
	Logger      profile.Logger
	Hystrix     profile.Hystrix
	Discovery   profile.Discovery
	Application Application
	Mysql       profile.Mysql
}

var configuration Profile

func Parse(f ...string) error {
	fileName := "app.toml"
	if len(f) > 0 {
		fileName = f[0]
	}

	return configuration.Init(fileName)
}

func initLogHook() error {
	os.Mkdir("log.d", 0666)
	path := "log.d/service.log"
	writer, err := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName("service.log"),       // 创建最新日志的链接
		rotatelogs.WithMaxAge(24*time.Hour),          // 每天一个日志文件
		rotatelogs.WithRotationTime(24*time.Hour*50), // 至多保存50天的日志
	)
	if err != nil {
		return fmt.Errorf("create log file error,%v", err)
	}

	logrus.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.ErrorLevel: writer,
		},
		&common.TextFormatter{}, // 自定义的日志格式
	))

	return nil
}

// InitLoggerWithProfile 初始化日志
func InitLoggerWithProfile(cfg *profile.Logger) error {
	switch cfg.Format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: cfg.TimeFormat,
		})
		logrus.Debug("Use json format logger")
	case "text", "":
		logrus.SetFormatter(&common.TextFormatter{
			UppercaseFirstMsgLetter: true,
		})
		logrus.Debug("Use text format logger")
	default:
		return fmt.Errorf("unsupport logrus formatter type %s", cfg.Format)
	}
	if cfg.Level != "" {
		logLevel, err := logrus.ParseLevel(cfg.Level)
		if err != nil {
			return fmt.Errorf("cannot parse logger level %s", cfg.Level)
		}
		logrus.SetLevel(logLevel)
	}

	if "" != cfg.LogFilePath {
		file, err := os.OpenFile(cfg.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		if nil != err {
			return fmt.Errorf("Open log file failed, err:%v, log file path:%v", err, cfg.LogFilePath)
		}
		logrus.SetOutput(file)
	}

	return initLogHook()
}

func makeStaticDiscover(lines []string) (func(string) ([]string, []string, error), error) {
	var staticServices []*discovery.StaticService
	for _, line := range lines {
		words := strings.Split(line, " ")
		if len(words) == 0 {

		}
		service := &discovery.StaticService{
			Name:  words[0],
			Hosts: words[1:],
		}

		staticServices = append(staticServices, service)
	}
	s := discovery.NewStaticDiscovery(staticServices)

	return s.Discover, nil
}

// Init 初始化配置
func (pro *Profile) Init(tomlFile string) error {
	_, _, err := profile.Parse(tomlFile, pro)
	if err != nil {
		return err
	}

	if err := InitLoggerWithProfile(&pro.Logger); err != nil {
		return err
	}

	// Discover 服务发现
	var discoverOption discovery.Option
	pro.Discovery.EnableConsul = false
	if pro.Discovery.EnableStatic {
		staticsDiscovery, err := makeStaticDiscover(pro.Discovery.StaticServices)
		if err != nil {
			return err
		}
		discoverOption.StaticFunc = staticsDiscovery
		logrus.Debug("static discovery enabled")
	}
	if err := discovery.Init(&discoverOption); err != nil {
		return err
	}

	// Invoker 服务调用初始化
	invokeOption := &invoke.Option{
		Discover:                     discovery.Discover,
		LoadBalanceMode:              pro.Invoker.LoadBanlanceMode,
		UseTracing:                   false,
		UseCircuit:                   pro.Invoker.CircuitEnabled,
		DoLogger:                     pro.Invoker.LoggerEnabled,
		DefaultErrorPercentThreshold: 20,
		DefaultTimeout:               18000,
	}

	return invoke.Init(invokeOption)
}

func Dump() {
	logrus.WithField("profile", jsonize.V(configuration, true)).Info("Dump profile")
}

// GetService 根据自己需要对方放出配置项目
// 返回服务配置
func GetService() *profile.Service {
	return &configuration.Service
}

func GetApplication() *Application {
	return &configuration.Application
}

func GetBase() *profile.Base {
	return &configuration.Base
}

func GetMysql() *profile.Mysql {
	return &configuration.Mysql
}
