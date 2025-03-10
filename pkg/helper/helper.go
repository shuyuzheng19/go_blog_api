package helper

import "fmt"

// CheckError 检查错误并处理
func CheckError(err error, message string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", message, err))
	}
}
