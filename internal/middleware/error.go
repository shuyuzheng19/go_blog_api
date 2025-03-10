package middleware

import (
	"blog/internal/handler"
	"blog/internal/utils"
	"blog/pkg/common"
	"blog/pkg/configs"
	"blog/pkg/logger"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func CustomRecover() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				path := c.Path()

				ip := utils.GetIPAddress(c)

				logError(ip, path, "未知错误", r)

				handler.ResultErrorToResponse(common.ERROR, c, "后端出现错误，请稍后重试.....")
			}
		}()
		return c.Next()
	}
}

// logError 记录错误日志
func logError(ip, path, errorType string, details interface{}) {
	// 使用三元运算符简化城市获取逻辑
	city := "开发环境"
	if configs.CONFIG.Env != "dev" {
		city = utils.GetIpCity(ip)
	}

	// 使用结构化日志记录错误信息
	logger.Error("记录错误日志",
		zap.String("path", path),
		zap.String("city", city),
		zap.String("type", errorType),
		zap.Any("error", details),
	)
}
