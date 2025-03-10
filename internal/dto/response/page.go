package response

// 分页返回
type Page struct {
	Page  int         `json:"page"`  //页码
	Count int64       `json:"total"` //总共多少条
	Size  int         `json:"size"`  //每页大小
	Data  interface{} `json:"data"`  //分页数据
}
