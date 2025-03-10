package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// EncryptPassword 对密码进行加密，使用 bcrypt 算法
func EncryptPassword(password string) string {
	// 生成一个随机的盐值并加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "" // 返回错误
	}
	return string(hashedPassword) // 返回加密后的密码
}

// VerifyPassword 验证输入的密码是否与存储的哈希密码匹配
func VerifyPassword(hashedPassword, password string) bool {
	// 使用 bcrypt.CompareHashAndPassword 函数比较密码
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println("密码验证失败:", err) // 记录错误
		return false                // 密码不匹配
	}
	return true // 密码匹配
}
