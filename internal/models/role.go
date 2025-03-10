package models

// Role 角色模型
type Role struct {
	Model
	ID          uint   `gorm:"primary_key;type:int;comment:角色ID" json:"id"`
	Name        string `gorm:"size:255;unique;not null;comment:角色名" json:"name"`
	Description string `gorm:"size:255;not null;comment:角色描述" json:"desc"`
}

func (*Role) TableName() string { return RoleTable }
