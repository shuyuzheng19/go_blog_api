package handler

import (
	"blog/internal/dto/requests"
	"blog/internal/dto/response"
	"blog/internal/models"
	"blog/internal/service"
	"blog/internal/utils"
	"blog/pkg/common"
	"blog/pkg/configs"

	"github.com/gofiber/fiber/v3"
)

func GetStatistics(ctx fiber.Ctx) error {
	var res, _ = service.GetSystemStatistics()
	return ResultSuccessToResponse(res, ctx)
}

func DeleteSystemInfoLog(ctx fiber.Ctx) error {
	var ids []int64

	if err := ctx.Bind().Body(&ids); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "参数绑定失败")
	}

	configs.DeleteData(models.SystemLogTable, nil, ids)

	return ResultSuccessToResponse(nil, ctx)
}

func GeSystemInfo(ctx fiber.Ctx) error {
	var res, _ = service.GetSystemInfo()
	return ResultSuccessToResponse(res, ctx)
}

func GetSystemLogInfo(ctx fiber.Ctx) error {
	// 1. 参数绑定
	var req requests.LogQueryParams
	if err := ctx.Bind().Query(&req); err != nil { // 使用 QueryParser 因为是 GET 请求
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "参数绑定失败")
	}

	// 2. 参数验证
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 3. 初始化查询构建器
	db := configs.DB.Model(&models.SystemLogInfo{})

	// 4. 构建查询条件
	if req.Keyword != nil && *req.Keyword != "" {
		like := "%" + *req.Keyword + "%"
		db = db.Where("message LIKE ? OR operator_name LIKE ? OR ip LIKE ? OR location LIKE ?", like, like, like, like)
	}

	if req.Module != nil && *req.Module != "" {
		db = db.Where("module = ?", *req.Module)
	}

	if req.Action != nil && *req.Action != "" {
		db = db.Where("action = ?", *req.Action)
	}

	if req.StartDate != nil && req.EndDate != nil {
		start, err1 := utils.ParseDateToTimestamp(*req.StartDate)
		end, err2 := utils.ParseDateToTimestamp(*req.EndDate)
		if err1 == nil && err2 == nil {
			db = db.Where("created_at BETWEEN ? AND ?", start, end)
		}
	}

	// 5. 执行查询
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "未找到相关数据")
	}

	var logs = make([]models.SystemLogInfo, 0)
	if total > 0 {
		offset := (req.Page - 1) * req.PageSize
		if err := db.Offset(offset).
			Limit(req.PageSize).
			Order(req.Sort.GetBlogOrderString("")). // 默认按创建时间倒序
			Find(&logs).Error; err != nil {
			return ResultErrorToResponse(common.ERROR, ctx, "查询数据失败")
		}
	}

	// 6. 返回结果
	return ResultSuccessToResponse(response.Page{
		Page:  req.Page,
		Size:  req.PageSize,
		Count: total,
		Data:  logs,
	}, ctx)
}
func GetSystemLogInfoLimit(ctx fiber.Ctx) error {
	var infos = make([]models.SystemLogInfo, 0)
	configs.DB.Model(&models.SystemLogInfo{}).Limit(10).
		Order(requests.CREATE.GetBlogOrderString("")).Find(&infos)
	return ResultSuccessToResponse(infos, ctx)
}
