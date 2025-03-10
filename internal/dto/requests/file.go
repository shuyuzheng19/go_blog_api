package requests

// FileRequest 查询过滤文件条件
// @Description 查询过滤文件条件
type FileRequest struct {
	Page    int    `form:"page"`    //第几页文件
	Keyword string `form:"keyword"` //文件的关键字
	Sort    string `form:"sort"`    //文件排序方式 size:大小排序和date:日期排序
}

type FileUpdateRequest struct {
	ID       int     `json:"id" validate:"required" error:"文件ID未传入"`
	Name     *string `json:"name"`
	IsPublic bool    `json:"is_pub"`
}

type TarRequest struct {
	Path string `json:"path" validate:"required" error:"压缩路径为空"`
	Min  int    `json:"min" validate:"required" error:"请填入定时删除分钟数"`
}

func (sort Sort) GetFilegOrderString(prefix string) string {
	switch sort {
	case CREATE:
		return prefix + "created_at desc"
	case UPDATE:
		return prefix + "updated_at  desc"
	case BACK:
		return prefix + "created_at asc"
	case SIZE:
		return prefix + "size desc"
	case BSIZE:
		return prefix + "size asc"
	default:
		return prefix + "created_at desc"
	}
}

// SystemFileRequest 过滤本地文件
type SystemFileRequest struct {
	Path    string `json:"path"`    //文件路径
	Keyword string `json:"keyword"` //文件关键字
}
