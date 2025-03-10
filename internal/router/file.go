package router

import (
	"blog/internal/handler"
	"blog/internal/middleware"
	"blog/pkg/common"

	"github.com/gofiber/fiber/v3"
)

// RegisterFileRouter 注册文件相关路由
func RegisterFileRouter(router fiber.Router) {
	fileController := handler.NewFileController()

	fileRouter := router.Group("/file")

	// 普通路由
	{
		fileRouter.Post("/avatar", fileController.UploadAvatar, middleware.LoggerMiddleware)

		fileRouter.Get("/public_list", fileController.GetPublicFileList)
	}

	// 管理员路由
	{
		fileRouter.Post("/upload", fileController.UploadFile, middleware.LoggerMiddleware, middleware.JwtMiddle(common.AdminRoleId), middleware.SystemLogMiddleware("file", "upload", "上传文件", false))

		fileRouter.Post("/upload/image", fileController.UploadImage, middleware.LoggerMiddleware, middleware.JwtMiddle(common.AdminRoleId), middleware.SystemLogMiddleware("file", "upload", "上传图片", false))

		fileRouter.Get("/current_list", fileController.GetCurrentFileFileList, middleware.JwtMiddle(common.AdminRoleId))

		fileRouter.Get("/admin/list", fileController.GetAdminFileList, middleware.AdminRequestMiddleware, middleware.JwtMiddle(common.AdminRoleId))

		fileRouter.Put("/admin/update", fileController.UpdateFileInfo, middleware.LoggerMiddleware, middleware.JwtMiddle(common.AdminRoleId))

		fileRouter.Put("/admin/delete", fileController.DeleteByIDs, middleware.LoggerMiddleware, middleware.JwtMiddle(common.AdminRoleId))
	}

	//超级管理员路由

	{
		fileRouter.Get("/admin/delete_md5", fileController.DeleteFileByMd5, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId))

		fileRouter.Get("/admin/get_upload_config", fileController.GetUploadConfig, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId))

		fileRouter.Post("/admin/set_upload_config", fileController.SetUploadConfig, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId))

		fileRouter.Post("/admin/system_file", fileController.GetSystemFile, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId))

		fileRouter.Get("/admin/system_file/clear_content", fileController.ClearSystemFileContent, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId))

		fileRouter.Post("/admin/system_file/delete", fileController.DeleteSystemFile, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId))

		fileRouter.Get("/admin/system_file/logs", fileController.GetLogFileList, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId))

		fileRouter.Get("/admin/system_file/current_log", fileController.GetCurrentLog, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId))

		fileRouter.Post("/admin/system_file/tar", fileController.TarDockerComposeData, middleware.LoggerMiddleware, middleware.JwtMiddle(common.SuperAdminRoleId))

		fileRouter.Get("/admin/system_file/tar", fileController.DownloadTar)

		fileRouter.Get("/admin/system_file/download", fileController.DownloadSystemFile)
	}
}
