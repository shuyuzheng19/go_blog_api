package handler

import (
	"blog/internal/dto/requests"
	"blog/internal/dto/response"
	"blog/internal/models"
	"blog/internal/service"
	"blog/internal/utils"
	"blog/pkg/common"
	"blog/pkg/logger"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type BlogController struct {
	service *service.BlogService
}

// CreateBlog 添加博客
func (b *BlogController) CreateBlog(ctx fiber.Ctx) error {
	var request requests.BlogRequest

	if err := ctx.Bind().Body(&request); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "无法解析请求体，请检查输入格式")
	}

	if errs := Validate(&request); len(errs) > 0 {
		return ResultValidatorErrorToResponse(ctx, errs)
	}

	if request.CategoryID == nil && request.TopicID == nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "请至少选择一个分类或专题")
	}

	uid := ctx.Locals("uid").(int)
	blog, err := b.service.CreateBlog(uid, request)
	if err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "创建博客时发生错误，请稍后重试")
	}

	logger.Info("博客创建成功", zap.Int("user_id", uid), zap.String("title", blog.Title))
	return ResultSuccessToResponse(blog, ctx)
}

// UpdateBlog 修改博客
func (b *BlogController) UpdateBlog(ctx fiber.Ctx) error {
	bid, err := strconv.ParseInt(ctx.Params("bid"), 10, 64)
	if err != nil || bid <= 0 {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "无效的博客ID")
	}

	var request requests.BlogRequest

	if err := ctx.Bind().Body(&request); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "无法解析请求体，请检查输入格式")
	}

	if errs := Validate(&request); len(errs) > 0 {
		return ResultValidatorErrorToResponse(ctx, errs)
	}

	if request.CategoryID == nil && request.TopicID == nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "请至少选择一个分类或专题")
	}

	if !request.IsPrivate {
		request.Password = nil
	}

	uid := ctx.Locals("uid").(int)
	rid := ctx.Locals("rid").(uint)

	if _, err := b.service.UpdateBlog(bid, uid, rid == uint(common.SuperAdminRoleId), request); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "更新博客时发生错误，请稍后重试")
	}

	logger.Info("博客更新成功", zap.Int64("id", bid), zap.Int("user_id", uid))
	return ResultSuccessToResponse(nil, ctx)
}

// GetBlogByID 根据 ID 获取博客
func (b *BlogController) GetBlogByID(ctx fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id", ""))
	if err != nil || id <= 0 {
		return ResultErrorToResponse(common.ERROR, ctx, "无效的博客ID")
	}

	blog, err := b.service.GetBlogByID(int64(id))

	if err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "无法获取博客，请稍后重试")
	}

	if blog == nil {
		return ResultErrorToResponse(common.ERROR, ctx, "博客不存在")
	}

	// 处理私密博客访问
	if blog.IsPrivate {
		password := ctx.Query("password", "")
		if password != *blog.Password && !b.isAuthorized(ctx, blog) {
			return ResultErrorToResponse(common.LockBlog, ctx, "该博客是私密的，需要密码访问")
		}
	}

	logger.Info("获取博客", zap.Int64("id", blog.ID), zap.String("title", blog.Title))

	blog.EyeCount = b.service.GetBlogEyeCount(blog.EyeCount, blog.ID)
	return ResultSuccessToResponse(blog.ToBlogContentResponse(), ctx)
}

// isAuthorized 检查用户是否有权访问私密博客
func (b *BlogController) isAuthorized(ctx fiber.Ctx, blog *models.Blog) bool {
	token := strings.TrimSpace(strings.TrimPrefix(ctx.Get("Authorization"), "Bearer"))
	if token == "" {
		return false
	}

	uid, rid := utils.ParseTokenUserIdAndRoleId(token)
	return uid == blog.UserID || rid == int(common.SuperAdminRoleId)
}

type editContent struct {
	Content string `json:"content"`
}

func (b *BlogController) SaveEditBlog(ctx fiber.Ctx) error {

	var edit editContent

	if err := ctx.Bind().Body(&edit); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "博客保存内容绑定失败")
	}

	var uid = ctx.Locals("uid").(int)

	var err = b.service.SaveEditBlog(uid, edit.Content)

	if err != nil {
		return ResultErrorToResponse(common.FAIL, ctx, "博客保存失败")
	}

	return ResultSuccessToResponse(nil, ctx)
}

func (b *BlogController) GetSaveEditBlog(ctx fiber.Ctx) error {

	var uid = ctx.Locals("uid").(int)

	var content, err = b.service.GetSaveEditBlog(uid)

	if err != nil {
		return ResultErrorToResponse(common.FAIL, ctx, "获取保存博客失败")
	}

	return ResultSuccessToResponse(content, ctx)
}

// GetBlogByIDToAdmin 管理员获取包含敏感信息博客
func (b *BlogController) GetBlogByIDToAdmin(ctx fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("bid", "0"))
	if err != nil || id <= 0 {
		return ResultErrorToResponse(common.ERROR, ctx, "无效的博客ID")
	}

	uid := ctx.Locals("uid").(int)
	rid := ctx.Locals("rid").(uint)

	blog, err := b.service.GetBlogByID(int64(id))
	if err != nil || blog == nil {
		return ResultErrorToResponse(common.ERROR, ctx, "无法获取博客，请稍后重试")
	}

	if rid == uint(common.SuperAdminRoleId) || uid == blog.UserID {
		maps := map[string]interface{}{
			"blog":      blog.ToBlogContentResponse(),
			"isPrivate": blog.IsPrivate,
			"password":  blog.Password,
		}
		return ResultSuccessToResponse(maps, ctx)
	}

	return ResultErrorToResponse(common.LockBlog, ctx, "您没有权限查看此博客")
}

// GetBlogList 获取博客列表
func (b *BlogController) GetBlogList(ctx fiber.Ctx) error {

	prequest := ctx.Locals(common.PageRequest).(requests.RequestQuery)

	var page = response.Page{
		Page: prequest.Page,
		Size: prequest.Size,
	}

	if err := b.service.GetBlogList(prequest, &page); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "无法获取博客列表，请稍后重试")
	}

	return ResultSuccessToResponse(page, ctx)
}

// GetAllAdminBlogList 获取所有管理员博客列表
func (b *BlogController) GetAllAdminBlogList(ctx fiber.Ctx) error {
	prequest := ctx.Locals(common.AdminRequest).(requests.AdminFilterRequest)

	var page = response.Page{
		Page: prequest.Page,
		Size: prequest.Size,
	}

	if err := b.service.GetAdminBlogList(nil, prequest, &page); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "无法获取管理员博客列表，请稍后重试")
	}

	return ResultSuccessToResponse(page, ctx)
}

// GetCurrentUserAdminBlogList 获取当前用户的管理员博客列表
func (b *BlogController) GetCurrentUserAdminBlogList(ctx fiber.Ctx) error {
	prequest := ctx.Locals(common.AdminRequest).(requests.AdminFilterRequest)

	var page = response.Page{
		Page: prequest.Page,
		Size: prequest.Size,
	}

	uid := ctx.Locals("uid").(int)

	if err := b.service.GetAdminBlogList(&uid, prequest, &page); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "无法获取当前用户的博客列表，请稍后重试")
	}

	return ResultSuccessToResponse(page, ctx)
}

// GetBlogArchive 获取博客归档
func (b *BlogController) GetBlogArchive(ctx fiber.Ctx) error {
	var req requests.ArchiveBlogRequest

	if err := ctx.Bind().Query(&req); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "无法解析归档请求参数")
	}

	var page = response.Page{}
	if err := b.service.GetArchiveBlog(req, &page); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "无法获取博客归档，请稍后重试")
	}

	return ResultSuccessToResponse(page, ctx)
}

// DeleteByIds 删除博客
func (b *BlogController) DeleteByIds(ctx fiber.Ctx) error {
	var ids []int64

	if err := ctx.Bind().Body(&ids); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "无法解析请求体，请检查输入格式")
	}

	uid := ctx.Locals("uid").(int)
	var userId *int
	rid := ctx.Locals("rid").(uint)

	if rid != uint(common.SuperAdminRoleId) {
		userId = &uid
	}

	if err := b.service.DeleteBlogByIDs(userId, ids); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "删除博客时发生错误，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

func (b *BlogController) InitSearch(ctx fiber.Ctx) error {
	go b.service.InitSearch()
	return ResultSuccessToResponse(nil, ctx)
}

func (b *BlogController) InitEyeCount(ctx fiber.Ctx) error {
	go b.service.InitEyeCount()
	return ResultSuccessToResponse(nil, ctx)
}

// UnDeleteByIds 恢复删除的博客
func (b *BlogController) UnDeleteByIds(ctx fiber.Ctx) error {
	var ids []int64

	if err := ctx.Bind().Body(&ids); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "无法解析请求体，请检查输入格式")
	}

	uid := ctx.Locals("uid").(int)
	var userId *int
	rid := ctx.Locals("rid").(uint)

	if rid != uint(common.SuperAdminRoleId) {
		userId = &uid
	}

	if err := b.service.UnDeleteBlogByIDs(userId, ids); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "恢复博客时发生错误，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// SimilarBlog 获取相似博客
func (b BlogController) SimilarBlog(ctx fiber.Ctx) error {
	keyword := ctx.Query("keyword", "")
	if keyword == "" {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "搜索关键字不能为空")
	}

	list, err := b.service.SimilarBlog(keyword)
	if err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "获取相似博客时发生错误，请稍后重试")
	}

	return ResultSuccessToResponse(list, ctx)
}

func (b BlogController) GetIndexData(ctx fiber.Ctx) error {
	var pinned = b.service.GetPinnedBlog()

	var page = response.Page{
		Page: 1,
		Size: 10,
	}

	b.service.GetBlogList(requests.RequestQuery{Page: 1, Sort: requests.CREATE, Cid: nil, Size: 10}, &page)

	return ResultSuccessToResponse(response.HomeData{Pinned: pinned, Page: page}, ctx)
}

func (b BlogController) SetPinnedBlog(ctx fiber.Ctx) error {
	var req requests.PinnedBlogRequest

	if err := ctx.Bind().Body(&req); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "参数验证失败")
	}

	var err = b.service.SetPinnedBlog(req)

	if err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, err.Error())
	}

	logger.Info("修改置顶博客成功", zap.Bool("pinned", req.Pinned), zap.Int64("blog_id", req.Id), zap.Int64p("order", req.Order))

	return ResultSuccessToResponse(nil, ctx)
}

// SearchBlog 搜索博客
func (b BlogController) SearchBlog(ctx fiber.Ctx) error {
	var req requests.SearchBlogRequest

	if err := ctx.Bind().Query(&req); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "无法解析搜索请求参数")
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	result := b.service.SearchBlog(req)
	logger.Info("博客搜索成功", zap.String("keyword", req.Keyword), zap.Int("page", req.Page))
	return ResultSuccessToResponse(result, ctx)
}

func (b BlogController) SearchBlog2(ctx fiber.Ctx) error {

	var keyword = ctx.Query("keyword")

	if keyword == "" {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "关键字为空")
	}

	result := b.service.SearchBlog(requests.SearchBlogRequest{Keyword: keyword, Page: 1})

	return ResultSuccessToResponse(result.Data, ctx)
}

// GetRecommendBlog 获取推荐博客
func (b *BlogController) GetRecommendBlog(ctx fiber.Ctx) error {
	result, err := b.service.GetRecommend()
	if err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "获取推荐博客时发生错误，请稍后重试")
	}
	return ResultSuccessToResponse(result, ctx)
}

// GetHotBlog 获取热门博客
func (b *BlogController) GetHotBlog(ctx fiber.Ctx) error {
	result, err := b.service.GetHotBlogs()
	if err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "获取热门博客时发生错误，请稍后重试")
	}

	return ResultSuccessToResponse(result, ctx)
}

// GetLatestBlog 获取最新博客
func (b *BlogController) GetLatestBlog(ctx fiber.Ctx) error {
	result, err := b.service.GetLatestBlogs()
	if err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "获取最新博客时发生错误，请稍后重试")
	}

	return ResultSuccessToResponse(result, ctx)
}

// SetRecommendBlog 设置推荐博客
func (b *BlogController) SetRecommendBlog(ctx fiber.Ctx) error {
	var ids []int

	if err := ctx.Bind().Body(&ids); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "无法解析请求体，请检查输入格式")
	}

	if err := b.service.SaveRecommend(ids); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "设置推荐博客时发生错误，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

func (b *BlogController) SetTempBlog(ctx fiber.Ctx) error {
	var req requests.TmpBlog

	if err := ctx.Bind().Body(&req); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "参数错误")
	}

	if errs := Validate(req); len(errs) > 0 {
		return ResultValidatorErrorToResponse(ctx, errs)
	}

	var user = ctx.Locals("user").(*models.User)

	req.Create = time.Now().Unix()

	req.User = &response.SimpleUserResponse{
		ID:       user.ID,
		NickName: user.NickName,
	}

	id, err := b.service.SetTempBlog(req)

	if err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, err.Error())
	}

	return ResultSuccessToResponse(id, ctx)
}

func (b *BlogController) GetTempBlog(ctx fiber.Ctx) error {
	var id = ctx.Query("id")
	blog := b.service.GetTempBlog(id)
	if blog == nil {
		return ResultErrorToResponse(common.NOT_FOUND, ctx, "获取失败，可能已经失效了。。。")
	}
	return ResultSuccessToResponse(blog, ctx)
}

// NewBlogController 创建新的 BlogController 实例
func NewBlogController() *BlogController {
	return &BlogController{service: service.NewBlogService()}
}
