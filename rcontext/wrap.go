package rcontext

import (
	"github.com/gin-gonic/gin"
	"github.com/lworkltd/kits/service/restful/code"
)

// Wrap 封装一层，将上层业务与HTTP框架解耦
func Wrap(f func(ctx RContext) (interface{}, error)) func(ginContext *gin.Context) (interface{}, code.Error) {
	return func(ginContext *gin.Context) (interface{}, code.Error) {
		rContext, err := NewRContext(ginContext)
		if err != nil {
			return nil, convertCodeError(err)
		}
		d, err := f(rContext)

		return d, convertCodeError(err)
	}
}

func convertCodeError(err error) code.Error {
	if err == nil {
		return nil
	}

	cerr, ok := err.(code.Error)
	if !ok {
		return code.NewMcode("UNKOWN_ERROR", err.Error())
	}

	return cerr
}
