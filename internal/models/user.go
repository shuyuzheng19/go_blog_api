package models

import (
	"blog/internal/dto/response"
)

// User 用户表模型
type User struct {
	Model
	ID        int    `gorm:"primary_key;type:int;comment:用户ID" json:"id"`
	Username  string `gorm:"size:16;unique;not null;comment:用户账号" validate:"required,min=8,max=16" error:"用户名为必填项，且长度在8到16个字符之间" json:"username"`
	Password  string `gorm:"size:255;not null;comment:用户密码" validate:"required,min=8;max=16" error:"密码为必填项，且长度至少为8到16个字符之间" json:"password"`
	Email     string `gorm:"size:255;unique;not null;comment:用户邮箱" validate:"required,email" error:"邮箱为必填项，且格式不正确" json:"email"`
	Avatar    string `gorm:"default:'test.png';comment:用户头像" json:"avatar"`
	Status    bool   `gorm:"size:1;default:1;comment:账号状态" json:"status"`
	NickName  string `gorm:"size:50;not null;comment:用户名称" validate:"required,min=1,max=50" error:"昵称为必填项，且长度在1到50个字符之间" json:"nickname"`
	RoleID    uint   `gorm:"column:role_id;type:integer;comment:角色ID" json:"role_id"`
	LoginIP   string `gorm:"size:45;comment:登录IP" json:"login_ip"`
	LoginCity string `gorm:"size:50;comment:登录地点" json:"login_city"`
	RegIp     string `gorm:"size:45;comment:注册IP" json:"reg_ip"`
	RegCity   string `gorm:"size:50;comment:注册地点" json:"reg_city"`
	RegTime   int64  `gorm:"comment:注册时间" json:"reg_time"`
	LastLogin int64  `gorm:"comment:最后登录时间" json:"last_login"`
	Role      Role   `json:"role"` // Role 字段不需要验证
}

func (*User) TableName() string { return UserTable }

func (u *User) ToVo() response.UserResponse {

	if u.LoginIP == "" {
		u.LoginIP = "内网"
		u.LoginCity = "未知"
	}

	return response.UserResponse{
		ID:       u.ID,
		Nickname: u.NickName,
		Avatar:   u.Avatar,
		Role:     u.Role.Name,
		Username: u.Username,
		Ip:       u.LoginIP,
		City:     u.LoginCity,
	}
}

func (u *User) ToSimpleUser() response.SimpleUserResponse {
	return response.SimpleUserResponse{
		ID:       u.ID,
		NickName: u.NickName,
	}
}
