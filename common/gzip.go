package common

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"sync"
)

// By: Henry @2019-01-25
// gzip.Writer内部会分配大量的堆内存，若不进行池化，则gc消耗的cpu会远比压缩计算本身
// 更耗费CPU。

type gzipWriterObj struct {
	w *gzip.Writer
	b *bytes.Buffer
}

// gzipWriterPool gzip压缩用对象池
type gzipWriterPool struct {
	*sync.Pool
}

var gzipWritePool = newGzipWriterPool()

func newGzipWriterPool() *gzipWriterPool {
	return &gzipWriterPool{
		&sync.Pool{New: func() interface{} {
			var buf bytes.Buffer
			return &gzipWriterObj{
				w: gzip.NewWriter(&buf),
				b: &buf,
			}
		}},
	}
}

func (wp *gzipWriterPool) Get() *gzipWriterObj {
	return (wp.Pool.Get()).(*gzipWriterObj)
}

func (wp *gzipWriterPool) Put(gw *gzipWriterObj) {
	gw.b.Reset()
	gw.w.Reset(gw.b)
	wp.Pool.Put(gw)
}

// GzipData 将元数据进行压缩
func GzipData(body []byte) ([]byte, error) {
	gzipUtil := gzipWritePool.Get()
	defer gzipWritePool.Put(gzipUtil)

	buffer := gzipUtil.b
	w := gzipUtil.w

	w.Write(body)
	w.Flush()
	if err := w.Close(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// ungzipDataWithPool 将压缩的数据解压缩
func ungzipDataWithNoPool(body []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	d, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return d, err
}

// UngzipData 将压缩的数据解压缩
var UngzipData = ungzipDataWithNoPool
