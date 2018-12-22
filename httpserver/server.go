package httpserver

import (
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/lworkltd/kits/service/profile"
	"github.com/lworkltd/kits/service/restful/wrap"
)

var errPrefix string
var wrapper *wrap.Wrapper

// TODO: 在gin所监听的接口同时处理pprof
func initService(engine *gin.Engine, option *profile.Service) error {
	wrapper = wrap.New(&wrap.Option{
		Prefix: option.McodePrefix,
	})

	errPrefix = option.McodePrefix

	if option.PprofEnabled {
		if option.PathPrefix != "" {
			ginpprof.WrapGroup(engine.Group(option.PathPrefix))
		} else {
			ginpprof.Wrapper(engine)
		}
	}

	return nil
}

func Run(option *profile.Service) error {
	r := gin.New()
	r.Use(Cors())
	r.Use(gin.Recovery())
	//gin.SetMode(gin.ReleaseMode)
	if err := initService(r, option); err != nil {
		return err
	}

	root := r.Group("/")
	if option.PathPrefix != "" {
		root = root.Group(option.PathPrefix)
	}

	wrapper.Get(root, "/ping", Wrap(OK))
	wrapper.Post(root, "/ping", Wrap(OK))
	routeV1(root.Group("/userapi/v1"))

	httpServer := &tradeMux{
		ginEngine: r,
		addr:      option.Host,
	}

	return httpServer.run()
}

func routeV1(v1 *gin.RouterGroup) {
	wrapper.Post(v1, "/user/create", Wrap(CreateUserRequest))
	wrapper.Post(v1, "/user/update", Wrap(CreateUserRequest))
	wrapper.Post(v1, "/user", Wrap(QueryUserListRequest))
	wrapper.Post(v1, "/user/list", Wrap(QueryUserListRequest))
}
