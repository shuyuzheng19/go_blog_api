package middleware

import (
	"blog/internal/dto/requests"
	"blog/internal/handler"
	"blog/pkg/common"

	"github.com/gofiber/fiber/v3"
)

// PaginationMiddleware 解析分页参数的中间件
func PaginationMiddleware(ctx fiber.Ctx) error {
	var request requests.RequestQuery

	if err := ctx.Bind().Query(&request); err != nil {
		return handler.ResultErrorToResponse(common.BAD_REQUEST, ctx, "参数绑定失败")
	}

	if request.Page <= 0 {
		request.Page = 1
	}

	if request.Size < 10 {
		request.Size = 10
	} else if request.Size > 30 {
		request.Size = 30
	}

	// 将分页参数存储到上下文中
	ctx.Locals(common.PageRequest, request)

	// 继续执行后续中间件或处理程序
	return ctx.Next()
}
