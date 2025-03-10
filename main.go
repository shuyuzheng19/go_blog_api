package main

import (
	"blog/internal/job"
	"blog/internal/router"
	"blog/internal/utils"
	"blog/pkg/configs"
	"blog/pkg/logger"
	"blog/pkg/smail"
	"log"
	"os"
)

func init() {

	// 加载全局配置
	configs.LoadGlobalConfig("application.yml")

	// 初始化日志系统
	logger.InitLogger(configs.CONFIG.Logger)

	// 确保 LOGGER 已经初始化
	if logger.Logger == nil {
		log.Fatal("Logger is not initialized")
	}

	logger.Info("全局配置加载成功")
	logger.Info("日志初始化成功")

	// 初始化SMTP邮箱
	smail.InitSmtp(configs.CONFIG.Mail)
	logger.Info("SMTP邮箱初始化成功")

	if configs.CONFIG.Env != "dev" {
		// 初始化数据库
		log.Println("加载数据库")
		configs.LoadDBConfig(configs.CONFIG.Db)
		log.Println("数据库加载完毕")

	}

	log.Println("加载Redis")
	configs.LoadRedis(configs.CONFIG.Redis)
	log.Println("Redis加载完毕")

	// 加载IP数据库
	log.Println("加载IP数据库")
	utils.LoadIpDB(configs.CONFIG.IpDbPath)
	log.Println("IP数据库加载完毕")

	// 创建上传目录
	if err := os.MkdirAll(configs.CONFIG.Upload.Path, os.ModePerm); err != nil {
		log.Fatalf("创建上传目录失败: %v", err)
	}
	log.Println("上传目录创建完毕")

	// 加载搜索配置
	configs.LoadSearchConfig(configs.CONFIG.Search)
	log.Println("搜索配置加载完毕")
}

// 程序入口
func main() {
	// 获取服务器配置
	var config = &configs.CONFIG.Server

	// 创建新的服务器实例
	var server = router.NewServer(config)

	// 注册路由
	registerRouters(server)

	// 启动定时任务（如果启用）
	if configs.CONFIG.Server.Cron {
		job.StartJob()
	}

	// 启动服务器
	server.Run()
}

// registerRouters 注册所有路由
func registerRouters(server *router.Server) {
	server.AddRouter(router.RegisterUserRouter)
	server.AddRouter(router.RegisterFileRouter)
	server.AddRouter(router.RegisterBlogRouter)
	server.AddRouter(router.RegisterCategoryRouter)
	server.AddRouter(router.RegisterTagRouter)
	server.AddRouter(router.RegisterTopicRouter)
	server.AddRouter(router.RegisterDataBaseRouter)
	server.AddRouter(router.RegisterConsoleRouter)
}
