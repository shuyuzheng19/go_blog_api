package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 生成Token的函数
func GenerateToken(payload map[string]interface{}, secret string) (string, error) {
	// 创建一个新的token对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(payload))

	// 签名并返回完整的编码token
	return token.SignedString([]byte(secret))
}

// 解析Token的函数
func ParseToken(tokenString string, secret string) (jwt.MapClaims, error) {
	// 解析token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保token的签名方法是我们预期的
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// 返回token的有效载荷
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func main() {
	secret := "your_secret_key"
	payload := map[string]interface{}{
		"user_id": 123,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 过期时间
	}

	// 生成token
	token, err := GenerateToken(payload, secret)
	if err != nil {
		fmt.Println("Error generating token:", err)
		return
	}
	fmt.Println("Generated Token:", token)

	// 解析token
	claims, err := ParseToken(token, secret)
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return
	}
	fmt.Println("Parsed Claims:", claims)
}
