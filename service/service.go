package service

import (
	"time"

	"github.com/lvhuat/http-service-example/common"

	"github.com/lvhuat/http-service-example/dto"
	"github.com/lvhuat/http-service-example/pb/userpb"
	"github.com/lvhuat/http-service-example/rcontext"
	"github.com/lworkltd/kits/service/restful/code"
)

func CreateUser(rctx rcontext.RContext, req *userpb.CreateUserRequest) (interface{}, error) {
	hashedPwd, _ := common.Encrypt(common.HashKey, req.Password)
	user := dto.User{
		UserName:       req.UserName,
		Address:        req.Address,
		Birthday:       req.Birthday,
		Email:          req.Email,
		Mobile:         req.Mobile,
		HashedPassword: hashedPwd,
		CreateTime:     time.Now().UnixNano() / int64(time.Millisecond),
	}

	if err := user.Insert(); err != nil {
		if dto.IsDuplicated(err) {
			return nil, code.NewMcodef(common.Error_UAPI_ALREADY_EXISTS.String(), "user name already exists")
		}
		return nil, err
	}

	return &userpb.CreateUserResponse{
		UserId: user.UserId,
	}, nil
}

func UpdateUser(rctx rcontext.RContext, req *userpb.UpdateUserRequest) (interface{}, error) {
	user := dto.User{
		UserName: req.UserName,
		Address:  req.Address,
		Birthday: req.Birthday,
		Email:    req.Email,
		Mobile:   req.Mobile,
	}

	if err := user.Update(); err != nil {
		return nil, err
	}

	return nil, nil
}

func QueryUser(rctx rcontext.RContext, req *userpb.QueryUserRequest) (interface{}, error) {
	user := dto.User{
		UserName: req.UserName,
	}

	if err := user.Load(); err != nil {
		return nil, err
	}

	return &userpb.QueryUserResponse{
		UserId:   user.UserId,
		UserName: user.UserName,
		Mobile:   user.Mobile,
		Email:    user.Email,
		Birthday: user.Birthday,
		Address:  user.Address,
	}, nil
}

func QueryUserList(rctx rcontext.RContext, req *userpb.QueryUserListRequest) (interface{}, error) {
	users, err := (&dto.User{}).GetList(req.CreateTime, req.Limit, req.Direct)
	if err != nil {
		return nil, err
	}

	userList := make([]*userpb.UserListItem, 0, len(users))

	for _, user := range users {
		userListItem := &userpb.UserListItem{
			UserId:   user.UserId,
			UserName: user.UserName,
			Mobile:   user.Mobile,
			Email:    user.Email,
			Birthday: user.Birthday,
			Address:  user.Address,
		}
		userList = append(userList, userListItem)
	}

	return &userpb.QueryUserListResponse{
		Users: userList,
	}, nil
}
