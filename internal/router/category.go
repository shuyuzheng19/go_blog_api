package router

import (
	"blog/internal/handler"
	"blog/internal/middleware"
	"blog/pkg/common"

	"github.com/gofiber/fiber/v3"
)

// RegisterCategoryRouter 分类相关路由
func RegisterCategoryRouter(router fiber.Router) {
	categoryController := handler.NewCategoryController()

	categoryRouter := router.Group("/category")

	// 普通用户路由
	{
		// 获取分类列表
		categoryRouter.Get("/list", categoryController.GetCategoryList)
	}

	// 管理员路由
	{
		// 目前没有普通管理员路由
	}

	// 超级管理员路由
	{
		// 获取所有管理员分类列表
		categoryRouter.Get("/admin/list", categoryController.GetAllAdminCategoryList, middleware.AdminRequestMiddleware, middleware.JwtMiddle(common.AdminRoleId))

		// 保存分类
		categoryRouter.Post("/admin/save", categoryController.SaveCategory, middleware.LoggerMiddleware, middleware.JwtMiddle(common.AdminRoleId), middleware.SystemLogMiddleware("category", "create", "添加分类", true))

		// 更新分类
		categoryRouter.Put("/admin/update", categoryController.UpdateCategory, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId), middleware.SystemLogMiddleware("category", "update", "修改分类", true))

		// 删除分类
		categoryRouter.Put("/admin/delete", categoryController.DeleteByIds, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId), middleware.SystemLogMiddleware("category", "delete", "删除分类", true))

		// 恢复删除的分类
		categoryRouter.Put("/admin/un_delete", categoryController.UnDeleteByIds, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId), middleware.SystemLogMiddleware("category", "undelete", "恢复分类", true))
	}
}
