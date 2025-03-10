package models

// Category 分类模型
type Category struct {
	Model
	ID   int    `gorm:"primary_key;type:int;comment:分类ID" json:"id"`
	Name string `gorm:"size:255;unique;not null;comment:分类名称" json:"name"`
}

// TableName 返回与模型对应的数据库表名
func (*Category) TableName() string {
	return CategoryTable
}
