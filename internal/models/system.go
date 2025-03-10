package models

type SystemLogInfo struct {
	Model
	ID           int64  `json:"id" gorm:"primaryKey"`
	Module       string `json:"module" gorm:"comment:模块名称"` // 如：blog、user、file等
	Action       string `json:"action" gorm:"comment:操作类型"` // 如：create、update、delete等
	Message      string `json:"message" gorm:"comment:消息"`
	IP           string `json:"ip" gorm:"comment:操作IP"`
	Location     string `json:"location" gorm:"comment:操作地点"`  // 如：北京市-朝阳区
	UserAgent    string `json:"userAgent" gorm:"comment:用户代理"` // 浏览器信息
	RequestURL   string `json:"requestUrl" gorm:"comment:请求URL"`
	Method       string `json:"method" gorm:"comment:请求方法"` // GET、POST等
	Params       string `json:"params" gorm:"comment:请求参数"`
	OperatorID   int    `json:"operatorId" gorm:"comment:操作人ID"`
	OperatorName string `json:"operatorName" gorm:"comment:操作人名称"`
	Email        string `json:"email" gorm:"comment:操作人邮箱"`
}

func (*SystemLogInfo) TableName() string {
	return SystemLogTable
}
