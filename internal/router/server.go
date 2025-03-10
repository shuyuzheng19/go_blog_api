package router

import (
	"blog/internal/middleware"
	"blog/pkg/configs"
	"blog/pkg/logger"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"go.uber.org/zap"
)

// Server 结构体，包含服务器配置和 Fiber 应用实例
type Server struct {
	config configs.ServerConfig // 服务器配置
	app    *fiber.App           // Fiber 应用实例
	router fiber.Router         // 路由器
}

// Run 启动服务器并监听指定地址
func (s Server) Run() {
	s.router.Get("/ping", pingHandler)
	logger.Info("服务器已启动，监听到", zap.String("addr", s.config.Addr))
	if err := s.app.Listen(s.config.Addr); err != nil {
		logger.Error("服务器启动失败", zap.Error(err))
		log.Fatal(err) // 启动服务器
	}
}

// pingHandler 健康检查路由处理函数
func pingHandler(c fiber.Ctx) error {
	logger.Info("PING")
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "PONG",
		"version": "1.0",
		"status":  "healthy",
		"text":    "书宇博客后台状态正常",
	})
}

// AddRouter 动态添加路由
func (s *Server) AddRouter(route func(router fiber.Router)) {
	route(s.router) // 调用传入的路由函数，传入当前的路由器
}

// AddMiddleware 添加中间件
func (s *Server) AddMiddleware(middleware ...func(ctx *fiber.Ctx) error) {
	for _, m := range middleware {
		s.router.Use(m)
	}
}

// NewServer 创建新的服务器实例
func NewServer(config *configs.ServerConfig) *Server {
	if config == nil {
		log.Fatalln("服务器启动失败，配置为空") // 配置为空时，记录错误并退出
	}

	logger.Info("启动服务器", zap.String("name", config.Name))

	// 配置 Fiber 应用
	fiberConfig := fiber.Config{
		AppName:      config.Name,                                      // 应用名称
		ServerHeader: "shuyu",                                          // 服务器标识
		BodyLimit:    config.MaxSize * 1024 * 1024,                     // 请求体大小限制：10MB
		ReadTimeout:  time.Second * time.Duration(config.ReadTimeOut),  // 读取超时
		WriteTimeout: time.Second * time.Duration(config.WriteTimeOut), // 写入超时
	}

	// 创建 Fiber 应用实例
	app := fiber.New(fiberConfig)

	// 使用自定义恢复中间件
	app.Use(middleware.CustomRecover())

	// 配置 CORS
	setupCORS(app, config.Cors)

	// 返回新的 Server 实例
	return &Server{config: *config, app: app, router: app.Group(config.ApiPrefix)}
}

// setupCORS 配置跨域资源共享
func setupCORS(app *fiber.App, corsConfig configs.CorsConfig) {
	if corsConfig.Enable {
		var corsConf cors.Config
		if !corsConfig.AllOrigins {
			corsConf = cors.Config{
				AllowOrigins:     corsConfig.AllowOrigins,
				AllowHeaders:     corsConfig.AllowHeaders,
				AllowMethods:     corsConfig.AllowMethods,
				AllowCredentials: corsConfig.AllowCredentials,
			}
		} else {
			corsConf = cors.ConfigDefault
		}
		app.Use(cors.New(corsConf))
		logger.Info("CORS 配置已启用", zap.Any("config", corsConfig))
	}
}
