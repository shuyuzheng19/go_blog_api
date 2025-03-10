package response

// UserResponse 用户信息
// @Description 返回用户的概要信息
type UserResponse struct {
	ID       int    `json:"id"`       //用户ID
	Nickname string `json:"nickName"` //用户名
	Avatar   string `json:"icon"`     //用户头像
	Role     string `json:"role"`     //用户角色
	Username string `json:"username"` //用户账号
	Ip       string `json:"ip"`       //登录IP
	City     string `json:"city"`     //登录地点
}

// TokenResponse 登陆成功返回的token概要
// @Description 登陆成功返回的token概要
type TokenResponse struct {
	Token  string       `json:"token"`  //token
	Expire string       `json:"expire"` //过期时间
	Create string       `json:"create"` //创建时间
	User   UserResponse `json:"user"`   //用户信息
}

type SimpleUserResponse struct {
	ID       int    `json:"id"`
	NickName string `json:"nickName"`
}
