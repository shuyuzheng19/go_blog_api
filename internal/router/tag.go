package router

import (
	"blog/internal/handler"
	"blog/internal/middleware"
	"blog/pkg/common"

	"github.com/gofiber/fiber/v3"
)

// RegisterTagRouter 标签相关路由
func RegisterTagRouter(router fiber.Router) {
	tagController := handler.NewTagController()

	tagRouter := router.Group("/tag")

	// 普通路由
	{
		// 获取随机标签列表
		tagRouter.Get("/random", tagController.GetTagRandomList)

		// 获取标签列表
		tagRouter.Get("/list", tagController.GetTagList)

		// 获取标签下的博客列表，带分页
		tagRouter.Get("/list/blog", tagController.GetTagBlogList, middleware.PaginationMiddleware)

		// 获取指定标签信息
		tagRouter.Get("/get/:tid", tagController.GetTagInfo)
	}

	// 管理员路由
	{
		// 目前没有普通管理员路由
	}

	// 超级管理员路由
	{
		// 获取所有管理员标签列表
		tagRouter.Get("/admin/list", tagController.GetAllAdminTagList, middleware.AdminRequestMiddleware, middleware.JwtMiddle(common.AdminRoleId))

		// 保存标签
		tagRouter.Post("/admin/save", tagController.SaveTag, middleware.LoggerMiddleware, middleware.JwtMiddle(common.AdminRoleId), middleware.SystemLogMiddleware("tag", "create", "创建标签", false))

		// 更新标签
		tagRouter.Put("/admin/update", tagController.UpdateTag, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId), middleware.SystemLogMiddleware("tag", "update", "修改标签", false))

		// 删除标签
		tagRouter.Put("/admin/delete", tagController.DeleteByIds, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId), middleware.SystemLogMiddleware("tag", "update", "删除标签", false))

		// 恢复删除的标签
		tagRouter.Put("/admin/un_delete", tagController.UnDeleteByIds, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId), middleware.SystemLogMiddleware("tag", "update", "恢复标签", false))
	}
}
