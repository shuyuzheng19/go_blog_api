package router

import (
	"blog/internal/handler"
	"blog/internal/middleware"
	"blog/pkg/common"

	"github.com/gofiber/fiber/v3"
)

func RegisterConsoleRouter(router fiber.Router) {
	consoleRouter := router.Group("/system")

	{
		consoleRouter.Get("/info", handler.GetStatistics, middleware.LoggerMiddleware, middleware.JwtMiddle(common.AdminRoleId))

		consoleRouter.Get("", handler.GeSystemInfo, middleware.LoggerMiddleware, middleware.JwtMiddle(common.AdminRoleId))

		consoleRouter.Get("/system_log", handler.GetSystemLogInfoLimit, middleware.LoggerMiddleware, middleware.JwtMiddle(common.AdminRoleId))

		consoleRouter.Get("/admin/system_log", handler.GetSystemLogInfo, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId))

		consoleRouter.Put("/admin/system_log", handler.DeleteSystemInfoLog, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId))
	}

}
