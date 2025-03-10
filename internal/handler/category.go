package handler

import (
	"blog/internal/dto/requests"
	"blog/internal/dto/response"
	"blog/internal/service"
	"blog/pkg/common"

	"github.com/gofiber/fiber/v3"
)

// CategoryController 分类控制器
type CategoryController struct {
	service *service.CategoryService
}

// GetCategoryList 获取分类列表
func (c *CategoryController) GetCategoryList(ctx fiber.Ctx) error {
	result := c.service.GetCategoryList()
	return ResultSuccessToResponse(result, ctx)
}

// GetAllAdminCategoryList 获取管理员分类列表
func (c *CategoryController) GetAllAdminCategoryList(ctx fiber.Ctx) error {
	prequest := ctx.Locals(common.AdminRequest).(requests.AdminFilterRequest)

	page := response.Page{
		Page: prequest.Page,
		Size: prequest.Size,
	}

	if err := c.service.GetAdminCategoryList(prequest, &page); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "无法获取分类列表，请稍后重试")
	}

	return ResultSuccessToResponse(page, ctx)
}

// SaveCategory 保存分类
func (c *CategoryController) SaveCategory(ctx fiber.Ctx) error {
	name := ctx.Query("name")

	if name == "" {
		return ResultErrorToResponse(common.FAIL, ctx, "请提供分类名称")
	}

	if err := c.service.SaveCategory(name); err != nil {
		return ResultErrorToResponse(common.FAIL, ctx, "添加分类失败，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// UpdateCategory 更新分类
func (c *CategoryController) UpdateCategory(ctx fiber.Ctx) error {
	var req requests.CategoryRequest

	if err := ctx.Bind().Body(&req); err != nil {
		return ResultErrorToResponse(common.FAIL, ctx, "请求参数无效，请检查输入")
	}

	// 参数验证
	if req.ID <= 0 || req.Name == "" {
		return ResultErrorToResponse(common.FAIL, ctx, "分类ID和名称均不能为空")
	}

	if err := c.service.UpdateCategory(req); err != nil {
		return ResultErrorToResponse(common.FAIL, ctx, "更新分类失败，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// DeleteByIds 批量删除分类
func (c *CategoryController) DeleteByIds(ctx fiber.Ctx) error {
	var ids []int64

	if err := ctx.Bind().Body(&ids); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "请求参数无效，请检查输入")
	}

	// 参数验证
	if len(ids) == 0 {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "请提供要删除的分类ID列表")
	}

	if err := c.service.DeleteByIDs(ids); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "删除分类失败，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// UnDeleteByIds 批量恢复分类
func (c *CategoryController) UnDeleteByIds(ctx fiber.Ctx) error {
	var ids []int64

	if err := ctx.Bind().Body(&ids); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "请求参数无效，请检查输入")
	}

	// 参数验证
	if len(ids) == 0 {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "请提供要恢复的分类ID列表")
	}

	if err := c.service.UnDeleteByIDs(ids); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "恢复分类失败，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// NewCategoryController 创建分类控制器实例
func NewCategoryController() *CategoryController {
	return &CategoryController{
		service: service.NewCategoryService(),
	}
}
