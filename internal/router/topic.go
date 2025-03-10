package router

import (
	"blog/internal/handler"
	"blog/internal/middleware"
	"blog/pkg/common"

	"github.com/gofiber/fiber/v3"
)

// RegisterTopicRouter 专题相关路由
func RegisterTopicRouter(router fiber.Router) {
	topicController := handler.NewTopicController()

	topicRouter := router.Group("/topic")

	// 普通路由
	{
		// 获取专题列表
		topicRouter.Get("/list", topicController.GetTopicByPage, middleware.LoggerMiddleware)

		// 获取指定专题下的博客列表，带分页
		topicRouter.Get("/blog/list/:tid", topicController.GetTopicBlogList, middleware.PaginationMiddleware)

		// 获取指定专题信息
		topicRouter.Get("/get/:tid", topicController.GetTopicInfo)

		// 获取专题下的所有博客
		topicRouter.Get("/blogs/:tid", topicController.GetTopicBlogs)
	}

	// 管理员路由
	{
		// 获取所有专题列表
		topicRouter.Get("/admin/sim_list", topicController.GetAllTopicList, middleware.JwtMiddle(common.AdminRoleId))
	}

	// 超级管理员路由
	{
		// 获取所有管理员专题列表
		topicRouter.Get("/admin/list", topicController.GetAllAdminTopicList, middleware.AdminRequestMiddleware, middleware.JwtMiddle(common.AdminRoleId))

		// 保存专题
		topicRouter.Post("/admin/save", topicController.SaveTopic, middleware.LoggerMiddleware, middleware.JwtMiddle(common.AdminRoleId), middleware.SystemLogMiddleware("topic", "create", "创建专题", true))

		// 更新专题
		topicRouter.Put("/admin/update", topicController.UpdateTopic, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId), middleware.SystemLogMiddleware("topic", "update", "修改专题", true))

		// 删除专题
		topicRouter.Put("/admin/delete", topicController.DeleteByIds, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId), middleware.SystemLogMiddleware("topic", "delete", "删除专题", true))

		// 恢复删除的专题
		topicRouter.Put("/admin/un_delete", topicController.UnDeleteByIds, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId), middleware.SystemLogMiddleware("topic", "un_delete", "恢复专题", true))
	}
}
