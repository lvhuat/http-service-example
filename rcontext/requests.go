package rcontext

import (
	"strconv"
	"time"

	"github.com/lvhuat/http-service-example/pb/userpb"
)

func (ctx *rContext) CreateUserRequest() *userpb.CreateUserRequest {
	req := &userpb.CreateUserRequest{}
	ctx.mustReadJson(req)
	panicIf(ctx, req.UserName == "", "invalid 'userName',missing")
	panicIf(ctx, req.Password == "", "invalid 'password',missing")
	return req
}

func (ctx *rContext) UpdateUserRequest() *userpb.UpdateUserRequest {
	req := &userpb.UpdateUserRequest{}
	ctx.mustReadJson(req)
	panicIf(ctx, req.UserName == "", "invalid 'userName',missing")
	return req
}

func (ctx *rContext) QueryUserRequest() *userpb.QueryUserRequest {
	req := &userpb.QueryUserRequest{}
	req.UserName = ctx.ginContext.Query("userName")
	panicIf(ctx, req.UserName == "", "invalid 'userName'")
	return req
}

func (ctx *rContext) QueryUserListRequest() *userpb.QueryUserListRequest {
	req := &userpb.QueryUserListRequest{}
	strLimit := ctx.ginContext.Query("limit")
	if strLimit == "" {
		req.Limit = 200
	} else {
		n, err := strconv.Atoi(strLimit)
		if err != nil || n < 0 || n > 3000 {
			panicIf(ctx, true, "invalid 'limit'")
		}
		req.Limit = int32(n)
	}

	strCreateTime := ctx.ginContext.Query("createTime")
	if strCreateTime == "" {
		req.CreateTime = time.Now().UnixNano() / int64(time.Millisecond)
	} else {
		n, err := strconv.Atoi(strCreateTime)
		if err != nil || n < 0 {
			panicIf(ctx, true, "invalid 'createTime'")
		}
		req.CreateTime = int64(n)
	}

	req.Direct = ctx.ginContext.Query("direct")
	switch req.Direct {
	case "":
	case userpb.QueryDirect_NEXT.String():
	case userpb.QueryDirect_PREV.String():
	default:
		panicIf(ctx, true, "invalid 'direct'")
	}

	return req
}
