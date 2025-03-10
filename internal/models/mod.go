package models

import (
	"gorm.io/gorm"
)

type Model struct {
	CreatedAt int64          `gorm:"autoCreateTime;comment:创建时间" json:"created_at"` // 创建时间，自动生成
	UpdatedAt int64          `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"` // 修改时间，自动生成
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at"`          // 删除时间，自动生成
}

// 表名称
const (
	UserTable        = "users"
	RoleTable        = "roles"
	TagTable         = "tags"
	BlogTable        = "blogs"
	CategoryTable    = "categories"
	TopicTable       = "topics"
	FileInfoTable    = "file_infos"
	FileInfoMd5Table = "file_md5_infos"
	BlogTagTable     = "blogs_tags"
	EyeCountTable    = "eye_count"
	SystemLogTable   = "system_log_info"
	EditBlogTable    = "edit_blog"
)
