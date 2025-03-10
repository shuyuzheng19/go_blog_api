package models

// Tag 标签模型
type Tag struct {
	Model
	ID   int    `gorm:"primary_key;type:int;comment:标签ID;" json:"id"`
	Name string `gorm:"size:255;unique;not null;comment:标签名称" json:"name"`
}

// TableName 返回与模型对应的数据库表名
func (*Tag) TableName() string {
	return TagTable
}
