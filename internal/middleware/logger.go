package middleware

import (
	"blog/internal/models"
	"blog/internal/utils"
	"blog/pkg/configs"
	"blog/pkg/logger"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// LoggerMiddleware 记录对方IP等信息日志中间件
func LoggerMiddleware(c fiber.Ctx) error {
	start := time.Now() // 记录开始时间

	// 继续执行后续中间件和处理程序
	err := c.Next()

	latency := time.Since(start) // 计算延迟

	ip := utils.GetIPAddress(c)                                   // 获取客户端IP
	city := utils.GetIpCity(ip)                                   // 获取IP所在城市
	userAgent := utils.GetClientPlatformInfo(c.Get("User-Agent")) // 获取用户代理信息

	// 记录日志
	logger.Info("客户端信息",
		zap.String("ip", ip),
		zap.String("city", city),
		zap.String("user_agent", userAgent),
		zap.String("latency", fmt.Sprintf("%v", latency)),
		zap.String("path", c.OriginalURL()),
		zap.Int("status", c.Response().StatusCode()),
		zap.String("method", c.Method()),
	)

	return err // 返回处理结果
}
func SystemLogMiddleware(module string, action string, message string, body bool) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		ip, city := utils.GetIpAndCitp(ctx)
		user := ctx.Locals("user")

		logInfo := &models.SystemLogInfo{
			Module:     module,
			Action:     action,
			IP:         ip,
			Location:   city,
			RequestURL: ctx.OriginalURL(),
			Method:     ctx.Method(),
			Message:    message,
			UserAgent:  utils.GetClientPlatformInfo(ctx.Get("User-Agent")),
		}

		if body {
			logInfo.Params = string(ctx.Body())
		}

		if user != nil {
			if userInfo, ok := user.(*models.User); ok {
				logInfo.OperatorID = userInfo.ID
				logInfo.OperatorName = userInfo.NickName
				logInfo.Email = userInfo.Email
			} else {
				// Log or handle the case where user is not of type *models.User
				logger.Warn("User is not of type *models.User")
			}
		}

		if err := configs.DB.Model(&models.SystemLogInfo{}).Create(logInfo).Error; err != nil {
			// Log the error
			logger.Error("Failed to create system log", zap.Error(err))
		}

		return ctx.Next()
	}
}
