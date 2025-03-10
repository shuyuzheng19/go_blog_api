package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// ValidateEmail 验证邮箱格式是否有效
func ValidateEmail(email string) bool {
	// 使用常见的邮箱格式正则表达式
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	// 检查邮箱格式是否匹配
	if !re.MatchString(email) {
		return false
	}
	return true
}

// Serialize 将对象序列化为 JSON 字符串
func Serialize[T any](t T) string {
	buff, err := json.Marshal(t)
	if err != nil {
		// 返回空字符串作为默认值
		fmt.Println("序列化失败:", err)
		return ""
	}
	return string(buff)
}

// Deserialize 将 JSON 字符串反序列化为对象
func Deserialize[T any](str string) T {
	var result T
	err := json.Unmarshal([]byte(str), &result)
	if err != nil {
		// 返回类型的零值作为默认值
		fmt.Println("反序列化失败:", err)
		return result // result 是类型 T 的零值
	}
	return result
}

func IsImageFile(suffix string) bool {
	// 将文件名转换为小写字母，并获取文件扩展名
	suffix = strings.ToLower(suffix)

	// 检查文件扩展名是否为图片格式
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, allowedExt := range allowedExtensions {
		if suffix == allowedExt {
			return true
		}
	}

	return false
}

// // IsImageFile 检查给定的文件是否为图片文件
// func IsImageFile(file multipart.File) (bool, error) {
// 	// 读取文件的前 512 字节，用于检测 MIME 类型
// 	buffer := make([]byte, 512)
// 	if _, err := file.Read(buffer); err != nil {
// 		return false, err
// 	}
// 	// 重置文件指针到开头，以便后续操作
// 	if _, err := file.Seek(0, 0); err != nil {
// 		return false, err
// 	}
// 	// 使用 http.DetectContentType 检测 MIME 类型
// 	contentType := http.DetectContentType(buffer)
// 	// 检查 MIME 类型是否属于常见的图片格式
// 	return strings.HasPrefix(contentType, "image/"), nil
// }
