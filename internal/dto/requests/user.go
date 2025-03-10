package requests

import (
	"blog/internal/models"
	"blog/internal/utils"
	"time"
)

// UserRequest 用户注册请求体
// @Description 用户注册请求体
type UserRequest struct {
	Username string `json:"username" validate:"required,min=8,max=16" error:"账号要在8-16个字符之间"`      //用户账号
	Password string `json:"password" validate:"required,min=8,max=16" error:"密码要在8-16个字符之间"`      //用户密码
	Email    string `json:"email" validate:"required,email" error:"这不是正确的邮箱格式"`                   //用户邮箱
	NickName string `json:"nickName" validate:"required,max=50,min=1" error:"用户名称最低1个字符，不能超过50个"` //用户名称
	Code     string `json:"code" validate:"required,min=6,max=6" error:"验证码为6位数字"`                //邮箱验证码
}

// LoginRequest 账号登录请求体
// @Description 账号登录请求体
type LoginRequest struct { //账号登录请求体
	Username string `json:"username" validate:"required" error:"账号不能为空"` //账号
	Password string `json:"password" validate:"required" error:"密码不能为空"` //密码
}

// 转成model
func (r UserRequest) ToUserModel(ip string) models.User {
	var city = utils.GetIpCity(ip)
	return models.User{
		Username: r.Username,
		Password: utils.EncryptPassword(r.Password),
		Email:    r.Email,
		RegIp:    ip,
		RegCity:  city,
		RegTime:  time.Now().Unix(),
		NickName: r.NickName,
		RoleID:   1,
	}
}

// ContactRequest 联系我请求
// @Description 联系我请求模型
type ContactRequest struct {
	Name    string `json:"name"  validate:"required" error:"名字为必填项"`        //你的名字
	Email   string `json:"email" validate:"required,email" error:"错误的邮箱格式"` //你的邮箱
	Subject string `json:"subject" validate:"required" error:"主题是必填项"`      //邮件主题
	Content string `json:"content" validate:"required" error:"内容是必填项"`      //邮件内容
}

// 用户列表过滤
type UserAdminFilter struct {
	Page    int    `form:"page"`    //第几页
	Keyword string `form:"keyword"` //关键字
	Sort    Sort   `form:"sort"`    //排序方式
	Deleted bool   `form:"deleted"` //是否过滤删除
	Start   string `form:"date[0]"` //开始时间
	Pub     *bool  `form:"pub"`
	End     string `form:"date[1]"` //结束时间
}

// 修改密码
type ResetPassword struct {
	Password    string `json:"password"`
	Email       string `json:"email"`
	OldPassWord string `json:"old_password"`
}

// GetOrderString 博客列表排序方式
func (sort Sort) GetUserOrderString(prefix string) string {
	switch sort {
	case CREATE:
		return prefix + "created_at desc"
	case UPDATE:
		return prefix + "updated_at  desc"
	case BACK:
		return prefix + "created_at asc"
	case ID:
		return prefix + "id asc"
	default:
		return prefix + "created_at desc"
	}
}

// UpdateUserRequest 管理员修改用户信息请求结构体
type UpdateUserRequest struct {
	ID            *int    `json:"id" validate:"required" error:"ID为必填项"`
	Username      *string `json:"username" validate:"required,min=8,max=16" error:"用户名长度在8到16个字符之间"`
	Password      string  `json:"password"`
	Email         *string `json:"email" validate:"required,email" error:"邮箱格式不正确"`
	Avatar        *string `json:"avatar" validate:"required,url" error:"头像URL格式不正确"`
	NickName      *string `json:"nickname" validate:"required,min=2,max=50" error:"昵称长度在2到50个字符之间"`
	Status        *bool   `json:"status" validate:"required" error:"账号状态必须为布尔值"`
	RoleID        *uint   `json:"role_id" validate:"required" error:"角色ID必须为有效值"`
	CurrentRoleId uint    `json:"-"`
}

func (u *UpdateUserRequest) ToModel() models.User {
	if u.Status == nil {
		*u.Status = false
	}
	return models.User{
		ID:       *u.ID,
		Username: *u.Username,
		Password: utils.EncryptPassword(u.Password),
		Email:    *u.Email,
		Avatar:   *u.Avatar,
		NickName: *u.NickName,
		Status:   *u.Status,
		RoleID:   *u.RoleID,
	}
}
