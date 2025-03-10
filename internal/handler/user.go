package handler

import (
	"blog/internal/dto/dtos"
	"blog/internal/dto/requests"
	"blog/internal/dto/response"
	"blog/internal/service"
	"blog/internal/utils"
	"blog/pkg/common"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

// UserController 用户控制器，处理与用户相关的请求
type UserController struct {
	service *service.UserService // 用户服务
}

// SendCodeToEmail 发送验证码到邮箱
func (u *UserController) SendCodeToEmail(ctx fiber.Ctx) error {
	email := ctx.Query("email", "")

	if email == "" || !utils.ValidateEmail(email) {
		return ResultErrorToResponse(common.EmailValidate, ctx, "请输入有效的邮箱地址")
	}

	if err := u.service.SendCodeToEmail(email); err != nil {
		return ResultErrorToResponse(common.FAIL, ctx, "验证码发送失败，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// RegisteredUser 注册新用户
func (u *UserController) RegisteredUser(ctx fiber.Ctx) error {
	var userRequest requests.UserRequest

	if err := ctx.Bind().Body(&userRequest); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "注册信息无效，请检查输入")
	}

	if errs := Validate(&userRequest); len(errs) > 0 {
		return ResultValidatorErrorToResponse(ctx, errs)
	}

	ip := utils.GetIPAddress(ctx)

	user, err := u.service.RegisteredUser(ip, userRequest)
	if err != nil {
		return ResultErrorToResponse(common.FAIL, ctx, "注册失败，请稍后重试")
	}

	return ResultSuccessToResponse(user.ToVo(), ctx)
}

// Login 用户登录
func (u *UserController) Login(ctx fiber.Ctx) error {
	var loginRequest requests.LoginRequest

	if err := ctx.Bind().Body(&loginRequest); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "登录信息无效，请检查输入")
	}

	if errs := Validate(&loginRequest); len(errs) > 0 {
		return ResultValidatorErrorToResponse(ctx, errs)
	}

	tokenResponse, err := u.service.Login(loginRequest)
	if err != nil {
		return ResultErrorToResponse(common.LoginFail, ctx, "登录失败，请检查用户名和密码")
	}

	if err := ResultSuccessToResponse(tokenResponse, ctx); err != nil {
		return err
	}

	ip := utils.GetIPAddress(ctx)

	go u.updateUserStatus(tokenResponse.User.ID, ip)

	return nil
}

// updateUserStatus 更新用户登录状态
func (u *UserController) updateUserStatus(userID int, ip string) {
	request := dtos.UserLoginStatus{
		ID:        userID,
		LastLogin: time.Now().Unix(),
		LoginIp:   ip,
		LoginCity: utils.GetIpCity(ip),
	}
	u.service.UpdateUserStatus(request)
}

// GetUserInfo 获取当前用户信息
func (u *UserController) GetUserInfo(ctx fiber.Ctx) error {
	user := GetUserInfo(ctx)

	if user == nil {
		return ResultErrorToResponse(common.NoLogin, ctx, "未能获取用户信息，请重新登录")
	}

	return ResultSuccessToResponse(user.ToVo(), ctx)
}

func (u *UserController) ResetPassword(ctx fiber.Ctx) error {
	user := GetUserInfo(ctx)

	if user == nil {
		return ResultErrorToResponse(common.NoLogin, ctx, "未能获取用户信息，请重新登录")
	}

	var req requests.ResetPassword

	if err := ctx.Bind().Body(&req); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "参数绑定失败")
	}

	if req.Email != user.Email || !utils.VerifyPassword(user.Password, req.OldPassWord) {
		return ResultErrorToResponse(common.Unauthorized, ctx, "身份验证失败，不允许修改!")
	}

	if err := u.service.ResetPassword(user.ID, req.Password); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "密码修改失败")
	}

	return ResultSuccessToResponse(nil, ctx)
}

func (u *UserController) UpdateUserRole(ctx fiber.Ctx) error {
	uid, err := strconv.Atoi(ctx.Query("id", ""))
	if err != nil || uid <= 0 {
		return ResultErrorToResponse(common.NoLogin, ctx, "无效的用户ID")
	}

	rid, err := strconv.Atoi(ctx.Query("rid", "1"))
	if err != nil {
		return ResultErrorToResponse(common.NoLogin, ctx, "无效的角色ID")
	}

	if err := u.service.UpdateRoleID(uid, uint(rid)); err != nil {
		return ResultErrorToResponse(common.NoLogin, ctx, "角色更新失败，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// ContactMe 用户反馈
func (u *UserController) ContactMe(ctx fiber.Ctx) error {
	var req requests.ContactRequest

	if err := ctx.Bind().Body(&req); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "反馈信息无效，请检查输入")
	}

	if errs := Validate(&req); len(errs) > 0 {
		return ResultValidatorErrorToResponse(ctx, errs)
	}

	if err := u.service.Contact(req); err != nil {
		return ResultErrorToResponse(common.FAIL, ctx, "发送反馈失败，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// UpdateUser 更新用户信息
func (u *UserController) UpdateUser(ctx fiber.Ctx) error {
	var req requests.UpdateUserRequest

	if err := ctx.Bind().Body(&req); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "用户信息无效，请检查输入")
	}

	if errs := Validate(req); len(errs) > 0 {
		return ResultValidatorErrorToResponse(ctx, errs)
	}

	req.CurrentRoleId = GetUserInfo(ctx).RoleID

	user, err := u.service.UpdateUser(&req)
	if err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "用户信息更新失败，请稍后重试")
	}

	return ResultSuccessToResponse(user, ctx)
}

// GetAdminUserList 获取管理员用户列表
func (u *UserController) GetAdminUserList(ctx fiber.Ctx) error {
	prequest := ctx.Locals(common.PageRequest).(requests.RequestQuery)

	page := response.Page{Page: prequest.Page, Size: prequest.Size}

	var request requests.UserAdminFilter

	if err := ctx.Bind().Query(&request); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "请求参数无效，请检查输入")
	}

	if err := u.service.GetAdminUserList(request, &page); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "无法获取用户列表，请稍后重试")
	}

	return ResultSuccessToResponse(page, ctx)
}

// GetWebSiteConfig 获取网站配置
func (u *UserController) GetWebSiteConfig(ctx fiber.Ctx) error {
	return ResultSuccessToResponse(u.service.GetWebSiteConfig(), ctx)
}

// UpdateWebSiteConfig 修改网站配置
func (u *UserController) UpdateWebSiteConfig(ctx fiber.Ctx) error {
	var config response.BlogConfigInfo

	if err := ctx.Bind().Body(&config); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "配置参数无效，请检查输入")
	}

	if err := u.service.SetWebSiteConfig(config); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "网站配置更新失败，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// GetRedisKeys 获取所有 Redis 键
func (u *UserController) GetRedisKeys(ctx fiber.Ctx) error {
	keys, err := u.service.GetRedisKeys()
	if err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "无法获取 Redis 键，请稍后重试")
	}
	return ResultSuccessToResponse(keys, ctx)
}

// DelRedisKeys 删除 Redis 键
func (u *UserController) DelRedisKeys(ctx fiber.Ctx) error {
	var keys []string

	if err := ctx.Bind().Body(&keys); err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "请求参数无效，请检查输入")
	}

	if err := u.service.DeleteRedisKeys(keys); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "删除 Redis 键失败，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// MatchDelKeys 匹配删除 Redis 键
func (u *UserController) MatchDelKeys(ctx fiber.Ctx) error {
	key := ctx.Query("key", "")

	if key == "" {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "匹配删除的键不能为空")
	}

	if err := u.service.DeleteMatchKeys(key); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "匹配删除 Redis 键失败，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// Logout 用户登出
func (u *UserController) Logout(ctx fiber.Ctx) error {
	uid := ctx.Locals("uid")

	if uid == nil {
		return ResultErrorToResponse(common.NoLogin, ctx, "您未登录")
	}

	if err := u.service.Logout(uid.(int)); err != nil {
		return ResultErrorToResponse(common.ERROR, ctx, "登出失败，请稍后重试")
	}

	return ResultSuccessToResponse(nil, ctx)
}

// GetUserBlogList 获取用户博客列表
func (u *UserController) GetUserBlogList(ctx fiber.Ctx) error {
	prequest := ctx.Locals(common.PageRequest).(requests.RequestQuery)

	page := response.Page{
		Page: prequest.Page,
		Size: prequest.Size,
	}

	if prequest.Uid == nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "无效的用户ID")
	}

	err := u.service.GetUsetBlogByPage(prequest, &page)
	if err != nil {
		page.Data = make([]response.BlogResponse, 0)
	}

	return ResultSuccessToResponse(&page, ctx)
}

// GetUserBlogTop10 获取用户前10篇博客
func (u *UserController) GetUserBlogTop10(ctx fiber.Ctx) error {
	uid, err := strconv.Atoi(ctx.Params("uid"))
	if err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "无效的用户ID")
	}

	list, err2 := u.service.GetUsetBlogTop10(uid)
	if err2 != nil {
		list = make([]response.SimpleBlogResponse, 0)
	}

	return ResultSuccessToResponse(list, ctx)
}

// GetUserTopics 获取用户的主题列表
func (u *UserController) GetUserTopics(ctx fiber.Ctx) error {
	uid, err := strconv.Atoi(ctx.Params("uid"))
	if err != nil {
		return ResultErrorToResponse(common.BAD_REQUEST, ctx, "无效的用户ID")
	}

	list, err2 := u.service.GetUsetTopic(uid)
	if err2 != nil {
		list = make([]response.TopicResponse, 0)
	}

	return ResultSuccessToResponse(list, ctx)
}

// NewUserController 创建新的用户控制器
func NewUserController() *UserController {
	service := service.NewUserService()
	return &UserController{service: service}
}
