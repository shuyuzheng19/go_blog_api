package router

import (
	"blog/internal/handler"
	"blog/internal/middleware"
	"blog/pkg/common"

	"github.com/gofiber/fiber/v3"
)

// RegisterUserRouter 注册用户相关路由
func RegisterUserRouter(router fiber.Router) {
	userController := handler.NewUserController() // 创建用户控制器实例

	userRouter := router.Group("/user") // 创建用户路由组

	// 普通路由
	{
		// 发送验证码到邮箱
		userRouter.Get("/send_email", userController.SendCodeToEmail, middleware.LoggerMiddleware, middleware.SystemLogMiddleware("user", "email", "发送邮件", true))

		// 注册新用户
		userRouter.Post("/registered", userController.RegisteredUser, middleware.LoggerMiddleware, middleware.SystemLogMiddleware("user", "registered", "用户注册", true))

		// 用户登录
		userRouter.Post("/login", userController.Login, middleware.LoggerMiddleware, middleware.SystemLogMiddleware("user", "login", "用户登录", true))

		// 获取网站配置
		userRouter.Get("/config", userController.GetWebSiteConfig)

		// 用户反馈
		userRouter.Post("/contact_me", userController.ContactMe, middleware.LoggerMiddleware, middleware.SystemLogMiddleware("user", "contact", "联系我", true))

		userRouter.Get("/blog/list", userController.GetUserBlogList, middleware.PaginationMiddleware)

		userRouter.Get("/blog/top/:uid", userController.GetUserBlogTop10)

		userRouter.Get("/topics/:uid", userController.GetUserTopics)
	}

	// 用户路由
	{
		// 获取当前用户信息
		userRouter.Get("/auth/get", userController.GetUserInfo, middleware.JwtMiddle(common.UserRoleId))

		// 获取当前用户信息
		userRouter.Put("/auth/reset", userController.ResetPassword, middleware.JwtMiddle(common.UserRoleId), middleware.SystemLogMiddleware("user", "reset", "重置密码", true))

		// 用户登出
		userRouter.Get("/logout", userController.Logout, middleware.LoggerMiddleware, middleware.JwtMiddle(common.UserRoleId), middleware.SystemLogMiddleware("user", "logout", "退出登录", true))
	}

	// 管理员路由
	{
	}

	// 超级管理员路由
	{

		// 获取管理员用户列表
		userRouter.Get("/admin/list", userController.GetAdminUserList, middleware.PaginationMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId))

		// 更新用户信息
		userRouter.Put("/admin/update", userController.UpdateUser, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId), middleware.SystemLogMiddleware("user", "update", "修改用户", true))

		// 修改用户角色
		userRouter.Put("/admin/update_role", userController.UpdateUserRole, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId), middleware.SystemLogMiddleware("user", "update", "修改用户角色", true))

		// 获取所有 Redis 键
		userRouter.Get("/admin/redis_keys", userController.GetRedisKeys, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId))

		// 删除 Redis 键
		userRouter.Put("/admin/redis_keys", userController.DelRedisKeys, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId), middleware.SystemLogMiddleware("user", "delete", "删除Redis键", true))

		// 匹配删除 Redis 键
		userRouter.Delete("/admin/redis_match_delete", userController.MatchDelKeys, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId), middleware.SystemLogMiddleware("user", "delete", "匹配删除Redis键", true))

		// 修改网站配置
		userRouter.Put("/admin/config", userController.UpdateWebSiteConfig, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId), middleware.SystemLogMiddleware("user", "config", "修改网站配置", true))
	}
}
