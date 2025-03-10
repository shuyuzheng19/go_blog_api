package utils

import (
	"blog/internal/dto/response"
	"blog/internal/models"
	"blog/pkg/common"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 生成Token的函数
func GenerateToken(user models.User) (result response.TokenResponse, err error) {

	var create = time.Now()

	var expire = time.Now().Add(common.TokenExpire)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(map[string]interface{}{
		"user_id": user.ID,
		"role_id": user.RoleID,
		"exp":     expire,
	}))

	var genToken string

	genToken, err = token.SignedString([]byte(common.TokenEncrypted))

	if err != nil {
		return
	}

	result = response.TokenResponse{
		Token:  genToken,
		Expire: FormatDate(expire),
		Create: FormatDate(create),
	}

	return
}

func ParseTokenUserId(tokenString string) int {
	var maps, err = ParseToken(tokenString)

	if err != nil {
		return -1
	}

	return int(maps["user_id"].(float64))
}

func ParseTokenUserIdAndRoleId(tokenString string) (int, int) {
	var maps, err = ParseToken(tokenString)

	if err != nil {
		return -1, -1
	}

	return int(maps["user_id"].(float64)), int(maps["role_id"].(float64))
}

// 解析Token的函数
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	// 解析token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保token的签名方法是我们预期的
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(common.TokenEncrypted), nil
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
