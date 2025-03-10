package response

// BlogConfigInfo 个人相关配置信息
// @Description 网站配置
type BlogConfigInfo struct {
	Name         string   `json:"name"`         //用户名称
	Avatar       string   `json:"avatar"`       //用户头像
	Descriptions []string `json:"descriptions"` //网站描述
	Content      string   `json:"content"`      //公告内容
}

// GetDefaultBlogConfigInfo 默认网站配置
func GetDefaultBlogConfigInfo() BlogConfigInfo {
	return BlogConfigInfo{
		Name:         "",
		Avatar:       "https://www.shuyuz.com/static/291ad9ffa5254e02e7f618f5f51eed51.png",
		Descriptions: []string{"后端程序员一枚", "有好的需求可以联系作者"},
		Content:      "",
	}
}
