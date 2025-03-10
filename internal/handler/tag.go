package handler

import (
	"blog/internal/dto/requests"
	"blog/internal/dto/response"
	"blog/internal/service"
	"blog/pkg/common"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// TagController 标签控制器
type TagController struct {
	service *service.TagService
}

// GetTagRandomList 获取随机标签列表
func (c *TagController) GetTagRandomList(ctx fiber.Ctx) error {
	result := c.service.RandomTags()
	return ResultSuccessToResponse(result, ctx)
}

// GetTagList 获取标签列表
func (c *TagController) GetTagList(ctx fiber.Ctx) error {
	result := c.service.GetTagList()
	return ResultSuccessToResponse(result, ctx)
}

// GetTagInfo 获取标签信息
func (c *TagController) GetTagInfo(ctx fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("tid", ""))
	if err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "无效的标签ID")
	}

	result := c.service.GetTagByID(id)
	if result == nil {
		return ResultErrorToResponse(common.NOT_FOUND, ctx, "未找到该标签")
	}

	return ResultSuccessToResponse(result, ctx)
}

// GetAllAdminTagList 获取管理员标签列表
func (c *TagController) GetAllAdminTagList(ctx fiber.Ctx) error {
	prequest := ctx.Locals(common.AdminRequest).(requests.AdminFilterRequest)

	page := response.Page{
		Page: prequest.Page,
		Size: prequest.Size,
	}

	if err := c.service.GetAdminTagList(prequest, &page); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "获取标签列表失败")
	}

	return ResultSuccessToResponse(page, ctx)
}

// DeleteByIds 批量删除标签
func (c *TagController) DeleteByIds(ctx fiber.Ctx) error {
	var ids []int64
	if err := ctx.Bind().Body(&ids); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "参数格式错误")
	}

	if len(ids) == 0 {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "标签ID列表不能为空")
	}

	if err := c.service.DeleteByIDs(ids); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "删除标签失败")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// UnDeleteByIds 批量恢复标签
func (c *TagController) UnDeleteByIds(ctx fiber.Ctx) error {
	var ids []int64
	if err := ctx.Bind().Body(&ids); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "参数格式错误")
	}

	if len(ids) == 0 {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "标签ID列表不能为空")
	}

	if err := c.service.UnDeleteByIDs(ids); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "恢复标签失败")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// SaveTag 保存标签
func (c *TagController) SaveTag(ctx fiber.Ctx) error {
	name := ctx.Query("name")
	if name == "" {
		return ResultErrorToResponse(common.FAIL, ctx, "标签名称不能为空")
	}

	if err := c.service.SaveTag(name); err != nil {
		return ResultErrorToResponse(common.FAIL, ctx, "标签添加失败")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// UpdateTag 更新标签
func (c *TagController) UpdateTag(ctx fiber.Ctx) error {
	var req requests.TagRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ResultErrorToResponse(common.FAIL, ctx, "参数格式错误")
	}

	if req.ID <= 0 || req.Name == "" {
		return ResultErrorToResponse(common.FAIL, ctx, "标签ID和名称不能为空")
	}

	if err := c.service.UpdateTag(req); err != nil {
		return ResultErrorToResponse(common.FAIL, ctx, "标签修改失败")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// GetTagBlogList 获取标签相关的博客列表
func (c *TagController) GetTagBlogList(ctx fiber.Ctx) error {
	prequest := ctx.Locals(common.PageRequest).(requests.RequestQuery)

	if prequest.Tid == nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "标签ID不能为空")
	}

	page := response.Page{
		Page: prequest.Page,
		Size: prequest.Size,
	}

	if err := c.service.GetTagIdBlogs(prequest, &page); err != nil {
		page.Count = 0
		page.Data = make([]response.BlogResponse, 0)
	}

	return ResultSuccessToResponse(&page, ctx)
}

// NewTagController 创建标签控制器实例
func NewTagController() *TagController {
	return &TagController{service: service.NewTagService()}
}
