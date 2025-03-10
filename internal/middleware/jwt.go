package middleware

import (
	"blog/internal/handler"
	"blog/internal/models"
	"blog/internal/utils"
	"blog/pkg/common"
	"strings"

	"github.com/gofiber/fiber/v3"
)

const (
	tokenType   = "Bearer "       // token类型
	tokenHeader = "Authorization" // token请求头
)

// ParseToken 解析token并返回用户信息
func ParseToken(header string, c fiber.Ctx) *models.User {

	// 检查header是否为空或不以tokenType开头
	if header == "" || !strings.HasPrefix(header, tokenType) {
		// 返回未登录错误
		handler.ResultErrorToResponse(common.NoLogin, c, "你还未登录，请先登录")
		return nil
	}

	// 去除tokenType前缀，获取实际token
	token := strings.TrimPrefix(header, tokenType)

	// 解析token获取用户ID
	uid := utils.ParseTokenUserId(token)

	// 检查用户ID是否有效
	if uid == -1 {
		// 返回解析token失败错误
		handler.ResultErrorToResponse(common.ParseTokenError, c, "解析Token失败")
		return nil
	}

	// 从Redis获取token进行比对
	redisToken := common.GetToken(uid)

	// 检查token是否过期
	if redisToken != token {
		// 返回token过期错误
		handler.ResultErrorToResponse(common.TokenExpireError, c, "token可能已过期，请重新登录")
		return nil
	}

	// 获取用户信息
	user := common.GetJwtUser(uid)

	// 检查用户是否存在
	if user == nil {
		// 返回未授权错误
		handler.ResultErrorToResponse(common.Unauthorized, c, "认证失败，请登录")
		return nil
	}

	// 返回有效用户信息
	return user
}

// JwtMiddle 验证身份中间件
func JwtMiddle(roleId common.RoleId) fiber.Handler {

	return func(c fiber.Ctx) error {
		header := c.Get(tokenHeader) // 获取请求头中的token

		user := ParseToken(header, c) // 解析token并获取用户信息

		// 如果用户信息无效，直接返回错误
		if user == nil {
			return handler.ResultErrorToResponse(common.Forbidden, c, "你还未登录")
		}

		role := user.Role.ID // 获取用户角色ID

		// 判断用户权限
		if !isAuthorized(roleId, role) {
			return handler.ResultErrorToResponse(common.Forbidden, c, "你没有权限")
		}

		c.Locals("user", user) // 设置用户到上下文
		c.Locals("uid", user.ID)
		c.Locals("rid", role)
		return c.Next() // 继续执行后续中间件
	}
}

// 权限检查函数
func isAuthorized(roleId common.RoleId, role uint) bool {
	if role == uint(common.SuperAdminRoleId) {
		return true
	}
	// 检查用户角色是否有权限
	return roleId == common.UserRoleId || (roleId == common.AdminRoleId && role == uint(common.AdminRoleId))
}
