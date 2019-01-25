package common

import (
	"math/rand"
	"time"
)

var (
	// CharsetLattinAndNumbers 数字，拉丁字母小写，拉丁字母大写
	CharsetLattinAndNumbers = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var rander *rand.Rand

func init() {
	rander = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomString 随机
func RandomString(n int) string {
	return RandomStringWithCharset(n, CharsetLattinAndNumbers)
}

// RandomStringWithCharset 使用自定义字符生成随机串
func RandomStringWithCharset(n int, charset string) string {
	b := make([]byte, n)
	for index := 0; index < n; index++ {
		b[index] = charset[rander.Intn(len(charset))]
	}
	return string(b)
}

// RandomString64 生成长度位64得随机字符串
func RandomString64() string {
	return RandomString(64)
}
