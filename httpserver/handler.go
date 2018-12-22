package httpserver

import (
	"github.com/lvhuat/http-service-example/rcontext"
	"github.com/lvhuat/http-service-example/service"
)

func CreateUserRequest(rctx rcontext.RContext) (interface{}, error) {
	return service.CreateUser(rctx, rctx.CreateUserRequest())
}

func UpdateUserRequest(rctx rcontext.RContext) (interface{}, error) {
	return service.UpdateUser(rctx, rctx.UpdateUserRequest())
}

func QueryUserRequest(rctx rcontext.RContext) (interface{}, error) {
	return service.QueryUser(rctx, rctx.QueryUserRequest())
}

func QueryUserListRequest(rctx rcontext.RContext) (interface{}, error) {
	return service.QueryUserList(rctx, rctx.QueryUserListRequest())
}
