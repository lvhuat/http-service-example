package common

import (
	"sync"

	jsoniter "github.com/json-iterator/go"
)

// XFeignKey 服务间调用的业务上下文信息
type XFeignKey struct {
	// API追踪ID，由GATEWAY（或NGINX）分配，在微服务作为上下文内容传递，当打印接口处理日志打印时必须打印
	TraceId string `json:"traceId"`
}

const (
	// HeaderFeignKey FeignKey 头部
	HeaderFeignKey = "X-Feign-Key"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var feignKeyPool = sync.Pool{
	New: func() interface{} {
		return &XFeignKey{}
	},
}

// PutFeignKey 回收FeignKey
func PutFeignKey(feignKey *XFeignKey) {
	feignKeyPool.Put(feignKey)
}

func (feignKey *XFeignKey) reset() {
	feignKey.TraceId = ""
}

// NewFeignKeyFromJsonBytes 使用JSON反序列化为结构
func NewFeignKeyFromJsonBytes(bb []byte) (*XFeignKey, error) {
	v := feignKeyPool.New().(*XFeignKey)
	v.reset()
	if err := json.Unmarshal(bb, v); err != nil {
		return nil, err
	}
	return v, nil
}
