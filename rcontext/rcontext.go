package rcontext

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/lvhuat/http-service-example/common"
	"github.com/lvhuat/http-service-example/pb/userpb"
	"github.com/lworkltd/kits/service/restful/code"
)

// RContext 解决业务时获取请求参数的问题
type RContext interface {
	DebugApi() bool
	GinContext() *gin.Context
	FeignKey() *common.XFeignKey
	ApiDesc() string
	LogEntry() *logrus.Entry

	CreateUserRequest() *userpb.CreateUserRequest
	UpdateUserRequest() *userpb.UpdateUserRequest
	QueryUserRequest() *userpb.QueryUserRequest
	QueryUserListRequest() *userpb.QueryUserListRequest
}

func IntOrDefault(s string, def int) int {
	if s == "" {
		return def
	}

	t, err := strconv.Atoi(s)
	if err != nil {
		return def
	}

	return t
}

func IntSliceOrDefault(s string, defs []int) []int {
	if s == "" {
		return defs
	}
	words := strings.Split(s, ",")
	slice := make([]int, 0, len(words))
	for _, word := range words {
		n, err := strconv.Atoi(word)
		if err != nil {
			return defs
		}
		slice = append(slice, n)
	}

	return slice
}

type rContext struct {
	ginContext *gin.Context
	feignKey   *common.XFeignKey
	apiDesc    string
}

func (ctx *rContext) FeignKey() *common.XFeignKey {
	return ctx.feignKey
}

func (ctx *rContext) LogEntry() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"traceId": ctx.FeignKey().TraceId,
		"apiDesc": ctx.ApiDesc(),
	})
}

// GinContext 返回HTTP的gin框架上下文
func (ctx *rContext) GinContext() *gin.Context {
	return ctx.ginContext
}

func (ctx *rContext) DebugApi() bool {
	return ctx.ginContext.GetHeader("X-Api-Debug") != ""
}

func (ctx *rContext) ApiDesc() string {
	//return ctx.ginContext.Request.Method + ":" + ctx.ginContext.Request.RequestURI
	return ctx.apiDesc
}

// newFeignKeyFromHttpRequest 从HTTP请求中提取FeignKey信息
func newFeignKeyFromHttpRequest(ctx *gin.Context) (*common.XFeignKey, error) {
	v := ctx.GetHeader(common.HeaderFeignKey)
	if v == "" {
		logrus.WithFields(logrus.Fields{
			"path": ctx.Request.RequestURI,
		}).Error("X-Feign-Key not found in HTTP header")

		return nil, code.NewMcodef("PARAMETER_ERROR", "X-Feign-Key not found in HTTP header")
	}

	return common.NewFeignKeyFromJsonBytes([]byte(v))
}

// NewRContext 创建一个处理HTTP请求的上下文
func NewRContext(ctx *gin.Context) (RContext, error) {
	fk, err := newFeignKeyFromHttpRequest(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"path":  ctx.Request.RequestURI,
		}).Error("newFeignKeyFromHttpRequest error")
		return nil, err
	}

	rctx := &rContext{
		ginContext: ctx,
		feignKey:   fk,
	}

	if ctx != nil {
		rctx.apiDesc = ctx.Request.Method + ":" + ctx.Request.RequestURI
	}

	return rctx, nil
}

func (ctx *rContext) readRequestBody() ([]byte, error) {
	request := ctx.ginContext.Request
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, code.New(-1, "bad request body")
	}
	defer request.Body.Close()

	return body, nil
}

func (ctx *rContext) readRequestJson(i interface{}) error {
	body, err := ctx.readRequestBody()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, i); err != nil {
		logrus.WithFields(logrus.Fields{
			"apiDesc": ctx.ApiDesc(),
			"traceId": ctx.FeignKey().TraceId,
			"body":    string(body),
		}).Error("readRequestBody failed")
		return code.New(403, "body json body")
	}

	logrus.WithFields(logrus.Fields{
		"apiDesc":  ctx.ApiDesc(),
		"traceId":  ctx.FeignKey().TraceId,
		"jsonBody": string(body),
	}).Debugln("Debug request")

	return nil
}

func (ctx *rContext) mustReadJson(i interface{}) {
	err := ctx.readRequestJson(i)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"traceId": ctx.FeignKey().TraceId,
		}).Error("Bad jason payload")
		panic(code.NewMcodef("PARAMETER_ERROR", "bad json payload"))
	}
}

func checkFeignKeyAccess(rctx *rContext) {
	fk := rctx.FeignKey()
	if fk == nil {
		return
	}
	// XXX
}

func panicIf(ctx *rContext, b bool, format string, args ...interface{}) {
	if !b {
		return
	}
	errMsg := fmt.Sprintf(format, args...)
	err := code.NewMcode("PARAMETER_ERROR", errMsg)
	if err != nil {
		ctx.LogEntry().WithFields(logrus.Fields{
			"error": errMsg,
		}).Errorln("Request with bad parameter")
	}

	panic(err)
}
