package requests

// AdminBlogFilterRequest 后台管理博客过滤条件
// @Description 后台管理博客过滤条件
type AdminFilterRequest struct {
	Page     int     `form:"page"`     //第几页
	Size     int     `form:"size"`     //每页多少条数据
	Keyword  *string `form:"keyword"`  //要搜索的关键字
	Sort     Sort    `form:"sort"`     //排序方式
	Category *int    `form:"category"` //指定分类
	Topic    *int    `form:"topic"`    //指定专题
	Start    *int64
	End      *int64
	Pub      *bool `form:"pub"`
	Deleted  bool  `form:"deleted"`
}

type LogQueryParams struct {
	Page      int     `json:"page" form:"page"`           // 页码
	PageSize  int     `json:"pageSize" form:"pageSize"`   // 每页数量
	Module    *string `json:"module" form:"module"`       // 模块
	Action    *string `json:"action" form:"action"`       // 操作类型
	Keyword   *string `json:"keyword" form:"keyword"`     // 关键词
	StartDate *string `json:"startDate" form:"startDate"` // 开始日期
	EndDate   *string `json:"endDate" form:"endDate"`     // 结束日期
	Start     *int64  `json:"-" form:"-"`
	End       *int64  `json:"-" form:"-"`
	Sort      Sort    `json:"sort" form:"sort"`
}
