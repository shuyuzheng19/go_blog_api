package router

import (
	"blog/internal/handler"
	"blog/internal/middleware"
	"blog/pkg/common"

	"github.com/gofiber/fiber/v3"
)

// RegisterTagRouter 标签相关路由
func RegisterDataBaseRouter(router fiber.Router) {
	dbController := handler.NewDataBaseController()

	dbRouter := router.Group("/database").Use(middleware.JwtMiddle(common.SuperAdminRoleId), middleware.LoggerMiddleware)

	{
		dbRouter.Get("get", dbController.GetTableInsertSQL)
		dbRouter.Post("exec", dbController.ExecSQL)
	}

}
