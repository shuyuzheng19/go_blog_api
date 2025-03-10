package common

//全局返回
// @Description 全局返回
type R struct {
	Code    Code        `json:"code"`           //返回的状态码
	Message string      `json:"message"`        //返回的消息
	Data    interface{} `json:"data,omitempty"` //返回的数据
}

// @Description 状态码
type Code uint

const (
	SUCCESS          Code = 200         //成功
	ERROR            Code = 500         //服务器错误
	Unauthorized     Code = 401         //认证失败
	NOT_FOUND        Code = 404         //找不到相关信息
	BAD_REQUEST      Code = 400         //请求参数错误
	Forbidden        Code = 403         //没权限访问
	FAIL             Code = iota + 1000 //后台处理失败
	EmailValidate    Code = 1001        //错误的邮箱格式
	NoLogin          Code = 1002        //未登录
	ParseTokenError  Code = 1003        //解析Token失败
	TokenExpireError Code = 1004        //Token错误
	LoginFail        Code = 1005        //登录失败
	LockBlog         Code = 1007        //博客加锁
)

var resultMaps = map[Code]string{
	SUCCESS:          "处理成功",
	ERROR:            "服务器错误",
	FAIL:             "处理失败",
	BAD_REQUEST:      "参数错误",
	NOT_FOUND:        "没有找到相关信息",
	Forbidden:        "你没有权限访问",
	EmailValidate:    "邮箱格式错误",
	NoLogin:          "你还未登录，请先登录",
	ParseTokenError:  "解析Token错误",
	TokenExpireError: "Token可能已过期，请登录后重试",
	Unauthorized:     "认证失败",
	LoginFail:        "登录失败",
	LockBlog:         "博客加锁",
}

func (c Code) DoData(data interface{}) R {
	return R{Code: c, Message: resultMaps[c], Data: data}
}

func (c Code) Do() R {
	return R{Code: c, Message: resultMaps[c]}
}
