package rcontext

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/lvhuat/http-service-example/common"
	"github.com/lvhuat/http-service-example/pb/userpb"
	"github.com/lworkltd/kits/service/restful/code"
	"github.com/sirupsen/logrus"
)

var (
	defaultParameterErrorCode = "PARAMETER_ERROR"
)

// RContext RContext是隔离网络层和业务层的中间层
type RContext interface {
	DebugApi() bool              // 用的表较少，返回表示本请求时debug属性，业务层可以根据需要做对应处理
	FeignKey() *common.XFeignKey // 大家用的比较多所以加上了
	ApiDesc() string             // Api的描述，用简短的语句表达本请求，比如 GET:/v1/user 或 GRPC:GetUser
	LogEntry() *logrus.Entry     // 默认日志打印，里面会带本请求的默认参数，比如userId，traceId，API名称等

	// 接口列表
	CreateUserRequest() *userpb.CreateUserRequest
	UpdateUserRequest() *userpb.UpdateUserRequest
	QueryUserRequest() *userpb.QueryUserRequest
	QueryUserListRequest() *userpb.QueryUserListRequest
}

// rContext 上下文实现
type rContext struct {
	ginContext *gin.Context
	feignKey   *common.XFeignKey
	apiDesc    string
	traceId    string
}

func (ctx *rContext) FeignKey() *common.XFeignKey {
	return ctx.feignKey
}

func (ctx *rContext) getTraceId() string {
	if ctx.traceId != "" {
		return ctx.traceId
	}
	if ctx.FeignKey() != nil {
		ctx.traceId = ctx.FeignKey().TraceId
	}

	if ctx.traceId == "" {
		ctx.traceId = common.RandomString(15)
	}

	return ctx.traceId
}

func (ctx *rContext) LogEntry() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"traceId": ctx.getTraceId(),
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
	return ctx.apiDesc
}

// newFeignKeyFromHttpRequest 从HTTP请求中提取FeignKey信息
func newFeignKeyFromHttpRequest(ctx *gin.Context) (*common.XFeignKey, error) {
	v := ctx.GetHeader(common.HeaderFeignKey)
	if v == "" {
		logrus.WithFields(logrus.Fields{
			"path": ctx.Request.RequestURI,
		}).Error("X-Feign-Key not found in HTTP header")

		return nil, code.NewMcodef(defaultParameterErrorCode, "X-Feign-Key not found in HTTP header")
	}

	return common.NewFeignKeyFromJsonBytes([]byte(v))
}

// NewRContext 创建一个处理HTTP请求的上下文
func NewRContext(ctx *gin.Context) (RContext, error) {
	// fk, err := newFeignKeyFromHttpRequest(ctx)
	// if err != nil {
	// 	logrus.WithFields(logrus.Fields{
	// 		"error": err,
	// 		"path":  ctx.Request.RequestURI,
	// 	}).Error("newFeignKeyFromHttpRequest error")
	// 	return nil, err
	// }

	rctx := &rContext{
		ginContext: ctx,
		//feignKey:   fk,
	}

	if ctx != nil {
		rctx.apiDesc = ctx.Request.Method + ":" + ctx.Request.RequestURI
	}

	return rctx, nil
}

// readRequestBody 读取http的body
func (ctx *rContext) readRequestBody() ([]byte, error) {
	request := ctx.ginContext.Request
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, code.NewMcodef(defaultParameterErrorCode, "bad request body")
	}
	defer request.Body.Close()

	return body, nil
}

// readRequestJson 将body解析为json
func (ctx *rContext) readRequestJson(i interface{}) error {
	body, err := ctx.readRequestBody()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, i); err != nil {
		logrus.WithFields(logrus.Fields{
			"apiDesc": ctx.ApiDesc(),
			"traceId": ctx.getTraceId(),
			"body":    string(body),
		}).Error("readRequestBody failed")
		panic(code.NewMcodef(defaultParameterErrorCode, "bad json payload"))
	}

	logrus.WithFields(logrus.Fields{
		"apiDesc":  ctx.ApiDesc(),
		"traceId":  ctx.FeignKey().TraceId,
		"jsonBody": string(body),
	}).Traceln("Debug request")

	return nil
}

// mustReadJson 如果解析json失败，则报出错误码异常
func (ctx *rContext) mustReadJson(i interface{}) {
	err := ctx.readRequestJson(i)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"traceId": ctx.FeignKey().TraceId,
		}).Error("Bad jason payload")
		panic(code.NewMcodef(defaultParameterErrorCode, "bad json payload"))
	}
}

func checkFeignKeyAccess(rctx *rContext) {
	fk := rctx.FeignKey()
	if fk == nil {
		return
	}
}

// panicIf 异常抛出函数
func panicIf(ctx *rContext, b bool, format string, args ...interface{}) {
	if !b {
		return
	}
	errMsg := fmt.Sprintf(format, args...)
	err := code.NewMcode(defaultParameterErrorCode, errMsg)
	if err != nil {
		ctx.LogEntry().WithFields(logrus.Fields{
			"error": errMsg,
		}).Errorln("Request with bad parameter")
	}

	panic(err)
}
