package httpserver

import (
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/lvhuat/http-service-example/rcontext"
	"github.com/lworkltd/kits/service/httpsrv"
	"github.com/lworkltd/kits/service/profile"
)

var errPrefix string
var wrapper *httpsrv.Wrapper
var log = logrus.WithField("pkg", "httpserver")

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func initService(option *profile.Service) error {
	wrapper = httpsrv.New(&httpsrv.Option{
		Prefix: option.McodePrefix,
	})

	errPrefix = option.McodePrefix

	if option.PprofEnabled {
		if option.PprofPathPrefix != "" {
			wrapper.Group(option.PprofPathPrefix).HandlePprof()
		} else {
			wrapper.HandlePprof()
		}
	}

	return nil
}

// Run 启动HTTP服务
func Run(option *profile.Service) error {
	if err := initService(option); err != nil {
		return err
	}

	root := wrapper.Group("/")
	if option.PathPrefix != "" {
		root = root.Group(option.PathPrefix)
	}

	wrapper.Any("/ping", rcontext.Wrap(OK))
	wrapper.HandleStat()

	v1 := wrapper.Group("/userapi/v1")
	v1.Post("/user/create", rcontext.Wrap(CreateUserRequest))
	v1.Post("/user/update", rcontext.Wrap(UpdateUserRequest))
	v1.Get("/user", rcontext.Wrap(QueryUserListRequest))
	v1.Post("/user/list", rcontext.Wrap(QueryUserListRequest))
	v1.HandlePprof()
	v1.HandleStat()
	httpServer := &tradeMux{
		wrapper: wrapper,
		addr:    option.Host,
	}

	return httpServer.run()
}

func routeV1(v1 *gin.RouterGroup) {
}
