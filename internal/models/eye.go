package models

type EyeView struct {
	Model
	ID    string `gorm:"primary_key;type:varchar(10);"` // 主键设置为日期字符串
	Count int64  `gorm:"type:int"`
}

func (*EyeView) TableName() string {
	return "eye_count"
}
