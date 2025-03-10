package handler

import (
	"blog/internal/dto/requests"
	"blog/internal/dto/response"
	"blog/internal/service"
	"blog/pkg/common"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// TopicController 专题控制器
type TopicController struct {
	service *service.TopicService
}

// GetTopicByPage 分页获取专题列表
func (t *TopicController) GetTopicByPage(ctx fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	result := t.service.GetTopicByPage(page)
	return ResultSuccessToResponse(result, ctx)
}

// GetTopicBlogList 获取专题博客列表
func (t *TopicController) GetTopicBlogList(ctx fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("tid", ""))
	if err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "无效的专题ID，请提供正确的ID")
	}

	prequest := ctx.Locals(common.PageRequest).(requests.RequestQuery)
	prequest.Tid = &id

	page := response.Page{
		Page: prequest.Page,
		Size: prequest.Size,
	}

	if err := t.service.GetTopicBlogList(prequest, &page); err != nil {
		page.Count = 0
		page.Data = make([]response.BlogResponse, 0)
	}

	return ResultSuccessToResponse(page, ctx)
}

// GetAllAdminTopicList 获取管理员专题列表
func (t *TopicController) GetAllAdminTopicList(ctx fiber.Ctx) error {
	prequest := ctx.Locals(common.AdminRequest).(requests.AdminFilterRequest)

	page := response.Page{
		Page: prequest.Page,
		Size: prequest.Size,
	}

	if err := t.service.GetAdminTopicList(prequest, &page); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "无法获取专题列表，请稍后重试")
	}

	return ResultSuccessToResponse(page, ctx)
}

// SaveTopic 保存专题
func (t *TopicController) SaveTopic(ctx fiber.Ctx) error {
	var req requests.TopicRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "专题参数格式错误，请检查输入")
	}

	if errs := Validate(&req); len(errs) > 0 {
		return ResultValidatorErrorToResponse(ctx, errs)
	}

	uid := ctx.Locals("uid").(int)
	if err := t.service.SaveTopic(uid, req); err != nil {
		return ResultErrorToResponse(common.FAIL, ctx, "无法添加专题，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// UpdateTopic 更新专题
func (t *TopicController) UpdateTopic(ctx fiber.Ctx) error {
	var req requests.TopicRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "专题参数格式错误，请检查输入")
	}

	if errs := Validate(&req); len(errs) > 0 {
		return ResultValidatorErrorToResponse(ctx, errs)
	}

	if req.ID <= 0 {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "无效的专题ID，请提供正确的ID")
	}

	uid := ctx.Locals("uid").(int)
	if err := t.service.UpdateTopic(uid, req); err != nil {
		return ResultErrorToResponse(common.FAIL, ctx, "无法更新专题，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// GetAllTopicList 获取所有专题列表
func (t *TopicController) GetAllTopicList(ctx fiber.Ctx) error {
	list, err := t.service.GetAllTopicList()
	if err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "无法获取专题列表，请稍后重试")
	}
	return ResultSuccessToResponse(list, ctx)
}

// GetTopicBlogs 获取专题博客
func (t *TopicController) GetTopicBlogs(ctx fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("tid", ""))
	if err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "无效的专题ID，请提供正确的ID")
	}

	result := t.service.GetTopicBlogs(id)
	return ResultSuccessToResponse(result, ctx)
}

// GetTopicInfo 获取专题信息
func (t *TopicController) GetTopicInfo(ctx fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("tid", ""))
	if err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "无效的专题ID，请提供正确的ID")
	}

	result := t.service.GetTopicInfo(id)
	if result == nil {
		return ResultErrorToResponse(common.NOT_FOUND, ctx, "未找到该专题，请检查ID")
	}

	return ResultSuccessToResponse(result, ctx)
}

// DeleteByIds 批量删除专题
func (t *TopicController) DeleteByIds(ctx fiber.Ctx) error {
	var ids []int64
	if err := ctx.Bind().Body(&ids); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "参数格式错误，请检查输入")
	}

	if len(ids) == 0 {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "专题ID列表不能为空，请提供至少一个ID")
	}

	if err := t.service.DeleteByIDs(ids); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "无法删除专题，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// UnDeleteByIds 批量恢复专题
func (t *TopicController) UnDeleteByIds(ctx fiber.Ctx) error {
	var ids []int64
	if err := ctx.Bind().Body(&ids); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "参数格式错误，请检查输入")
	}

	if len(ids) == 0 {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "专题ID列表不能为空，请提供至少一个ID")
	}

	if err := t.service.UnDeleteByIDs(ids); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "无法恢复专题，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// NewTopicController 创建专题控制器实例
func NewTopicController() *TopicController {
	return &TopicController{service: service.NewTopicService()}
}
