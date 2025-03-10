package utils

import (
	"math/rand"
	"strconv"
)

// 随机生成6位数字
func RandomNumberCode() string {
	return strconv.Itoa(rand.Intn(900000) + 100000)
}
