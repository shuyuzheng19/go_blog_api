package middleware

import (
	"blog/internal/dto/requests"
	"blog/internal/handler"
	"blog/internal/utils"
	"blog/pkg/common"

	"github.com/gofiber/fiber/v3"
)

// PaginationMiddleware 解析分页参数的中间件
func AdminRequestMiddleware(ctx fiber.Ctx) error {
	var request requests.AdminFilterRequest

	if err := ctx.Bind().Query(&request); err != nil {
		return handler.ResultErrorToResponse(common.BAD_REQUEST, ctx, "参数绑定失败")
	}

	if request.Page <= 0 {
		request.Page = 1
	}

	if request.Size < 10 {
		request.Size = 10
	}

	var date0 = ctx.Query("date[0]")

	var date1 = ctx.Query("date[1]")

	if date0 != "" && date1 != "" {
		var start, err = utils.ParseDateToTimestamp(date0)
		var end, err2 = utils.ParseDateToTimestamp(date1)
		if err == nil && err2 == nil {
			request.Start = &start
			request.End = &end
		}
	}

	// 将分页参数存储到上下文中
	ctx.Locals(common.AdminRequest, request)

	// 继续执行后续中间件或处理程序
	return ctx.Next()
}
