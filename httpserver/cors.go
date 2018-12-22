package httpserver

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lvhuat/http-service-example/rcontext"
	"github.com/lworkltd/kits/service/context"
	"github.com/lworkltd/kits/service/restful/code"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		var headerKeys []string
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("Access-Control-Allow-Origin, Access-Control-Allow-Headers, %s", headerStr)
		} else {
			headerStr = "Access-Control-Allow-Originn, Access-Control-Allow-Headers"
		}

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Expose-Headers", "*")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}

func codeError(err error) code.Error {
	if err == nil {
		return nil
	}

	cerr, ok := err.(code.Error)
	if !ok {
		return code.New(20000000, err.Error())
	}

	return cerr
}

var (
	// CerrCheckSnowProtect 雪崩预警
	CerrCheckSnowProtect = code.NewMcode("SNOWSLIDE_DENIED", "Check Snow Protect")
	gCurTime             int64
	checkSnowMutex       sync.Mutex
	gCurCount            int32
)

func checkSnowSlide(showCount int32) error {
	timeNow := time.Now().Unix()
	checkSnowMutex.Lock()
	defer checkSnowMutex.Unlock()

	if timeNow > gCurTime {
		gCurTime = timeNow
		gCurCount = 1
		return nil
	}
	if gCurCount >= showCount {
		return CerrCheckSnowProtect
	}
	gCurCount++

	return nil
}

func Wrap(f func(ctx rcontext.RContext) (interface{}, error)) func(srvContext context.Context, ginContext *gin.Context) (interface{}, code.Error) {
	return func(srvContext context.Context, ginContext *gin.Context) (interface{}, code.Error) {
		if ginContext.Request.URL.Path == "/ping" {
			return nil, nil
		}

		err := checkSnowSlide(3000)
		if err != nil {
			return nil, codeError(err)
		}

		rContext, err := rcontext.NewRContext(ginContext)
		if err != nil {
			return nil, codeError(err)
		}

		d, err := f(rContext)

		return d, codeError(err)
	}
}

func OK(ctx rcontext.RContext) (interface{}, error) {
	return nil, nil
}
