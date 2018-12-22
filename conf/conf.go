package conf

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/lworkltd/kits/helper/consul"
	"github.com/lworkltd/kits/service/discovery"
	"github.com/lworkltd/kits/service/invoke"
	"github.com/lworkltd/kits/service/profile"
	"github.com/lworkltd/kits/utils/eval"
	"github.com/lworkltd/kits/utils/ipnet"
	"github.com/lworkltd/kits/utils/jsonize"
	"github.com/lworkltd/kits/utils/log"
)

type Profile struct {
	Base        profile.Base
	Consul      profile.Consul
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

func (pro *Profile) Init(tomlFile string) error {
	_, _, err := profile.Parse(tomlFile, pro)
	if err != nil {
		return err
	}

	if err := log.InitLoggerWithProfile(&pro.Logger); err != nil {
		return err
	}

	consulClient, err := consul.New(pro.Consul.Url)
	if err != nil {
		return err
	}
	consul.SetClient(consulClient)

	// 注册eval解析器
	eval.RegisterKeyValueExecutor("kv_of_consul", consulClient.KeyValue)
	eval.RegisterKeyValueExecutor("ip_of_interface", ipnet.Ipv4)

	// 将填充配置中使用了eval语法
	if err := eval.Complete(&pro); err != nil {
		return err
	}

	// Discover 服务发现
	var discoverOption discovery.Option
	discoverOption.RegisterFunc = consulClient.Register
	discoverOption.UnregisterFunc = consulClient.Unregister
	if pro.Discovery.EnableConsul {
		discoverOption.SearchFunc = consulClient.Discover
		logrus.Debug("consul discovery enabled")
	}
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
	if err := invoke.Init(invokeOption); err != nil {
		return err
	}

	return nil
}

func Dump() {
	mutiline := log.IsMultiLineFormat(configuration.Logger.Format)
	logrus.WithField("profile", jsonize.V(configuration, mutiline)).Info("Dump profile")
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
