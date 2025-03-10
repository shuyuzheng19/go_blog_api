package router

import (
	"blog/internal/handler"
	"blog/internal/middleware"
	"blog/pkg/common"

	"github.com/gofiber/fiber/v3"
)

// RegisterBlogRouter 注册博客相关路由
func RegisterBlogRouter(router fiber.Router) {
	blogController := handler.NewBlogController()
	blogRouter := router.Group("/blog")

	// 普通路由
	{
		// 获取博客
		blogRouter.Get("/get/:id", blogController.GetBlogByID, middleware.LoggerMiddleware)

		//获取置顶的博客
		blogRouter.Get("/index", blogController.GetIndexData)

		// 获取博客列表
		blogRouter.Get("/list", blogController.GetBlogList, middleware.PaginationMiddleware, middleware.LoggerMiddleware)

		// 获取推荐博客
		blogRouter.Get("/recommend", blogController.GetRecommendBlog)

		// 获取热门博客
		blogRouter.Get("/hots", blogController.GetHotBlog)

		// 获取最新博客
		blogRouter.Get("/latest", blogController.GetLatestBlog)

		// 获取博客归档
		blogRouter.Get("/archive", blogController.GetBlogArchive)

		// 搜索博客
		blogRouter.Get("/search", blogController.SearchBlog, middleware.LoggerMiddleware)

		// 搜索博客
		blogRouter.Get("/search2", blogController.SearchBlog2)

		// 获取相似博客
		blogRouter.Get("/similar", blogController.SimilarBlog)
	}

	// 管理员路由
	{
		// 创建博客
		blogRouter.Post("/admin/create", blogController.CreateBlog, middleware.LoggerMiddleware, middleware.JwtMiddle(common.AdminRoleId), middleware.SystemLogMiddleware("blog", "create", "创建博客", false))

		// 修改博客
		blogRouter.Put("/admin/update/:bid", blogController.UpdateBlog, middleware.LoggerMiddleware, middleware.JwtMiddle(common.AdminRoleId), middleware.SystemLogMiddleware("blog", "update", "修改博客", false))

		// 获取当前用户的博客列表
		blogRouter.Get("/admin/user/list", blogController.GetCurrentUserAdminBlogList, middleware.LoggerMiddleware, middleware.AdminRequestMiddleware, middleware.JwtMiddle(common.AdminRoleId))

		// 删除博客
		blogRouter.Put("/admin/delete", blogController.DeleteByIds, middleware.LoggerMiddleware, middleware.JwtMiddle(common.AdminRoleId), middleware.SystemLogMiddleware("blog", "delete", "删除博客", true))

		// 恢复删除的博客
		blogRouter.Put("/admin/un_delete", blogController.UnDeleteByIds, middleware.LoggerMiddleware, middleware.JwtMiddle(common.AdminRoleId), middleware.SystemLogMiddleware("blog", "un_delete", "恢复博客", true))

		// 管理员获取博客
		blogRouter.Get("/admin/update/get/:bid", blogController.GetBlogByIDToAdmin, middleware.LoggerMiddleware, middleware.JwtMiddle(common.AdminRoleId), middleware.SystemLogMiddleware("blog", "update", "获取修改博客", true))

		// 保存编辑的博客内容
		blogRouter.Post("/admin/save_edit", blogController.SaveEditBlog, middleware.LoggerMiddleware, middleware.JwtMiddle(common.AdminRoleId), middleware.SystemLogMiddleware("blog", "create", "保存编辑博客内容", true))

		// 获取保存编辑的博客内容
		blogRouter.Get("/admin/get_edit", blogController.GetSaveEditBlog, middleware.LoggerMiddleware, middleware.JwtMiddle(common.AdminRoleId), middleware.SystemLogMiddleware("blog", "get", "获取保存编辑博客内容", true))

		//保存临时博客
		blogRouter.Post("/admin/save_temp", blogController.SetTempBlog, middleware.LoggerMiddleware, middleware.JwtMiddle(common.AdminRoleId), middleware.SystemLogMiddleware("blog", "create", "保存临时博客内容", true))

		//获取保存的临时博客
		blogRouter.Get("/admin/get_temp", blogController.GetTempBlog, middleware.LoggerMiddleware, middleware.SystemLogMiddleware("blog", "get", "获取保存临时博客内容", true))
	}

	// 超级管理员路由
	{
		// 设置推荐博客
		blogRouter.Post("/admin/recommend", blogController.SetRecommendBlog, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId), middleware.SystemLogMiddleware("blog", "un_delete", "设置推荐", true))

		// 获取所有管理员的博客列表
		blogRouter.Get("/admin/list", blogController.GetAllAdminBlogList, middleware.LoggerMiddleware, middleware.AdminRequestMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId))

		// 修改博客置顶
		blogRouter.Post("/admin/pinned", blogController.SetPinnedBlog, middleware.LoggerMiddleware, middleware.AdminRequestMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId))

		// 初始化搜索
		blogRouter.Get("/admin/init_search", blogController.InitSearch, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId), middleware.SystemLogMiddleware("blog", "init", "初始化搜索", false))

		// 初始化浏览量
		blogRouter.Get("/admin/init_eye", blogController.InitEyeCount, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId), middleware.SystemLogMiddleware("blog", "init", "初始化浏览量", false))
	}
}
