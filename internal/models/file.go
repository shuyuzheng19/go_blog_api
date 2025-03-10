package models

// FileInfo 文件信息结构体
type FileInfo struct {
	Model
	ID          int         `gorm:"primary_key;type:int;comment:文件ID"`
	OldName     string      `gorm:"size:255;not null;comment:原文件名"`
	NewName     string      `gorm:"size:255;not null;comment:新文件名"`
	UserID      *int        `gorm:"column:user_id;type:int;comment:上传文件用户ID"`
	Suffix      string      `gorm:"size:10;comment:后缀"`
	Size        int64       `gorm:"comment:文件大小"`
	FileMd5     string      `gorm:"size:255;column:md5;not null;comment:文件 MD5 值"`
	IsPub       bool        `gorm:"default:false;comment:是否公开"`
	FileMd5Info FileMd5Info `gorm:"foreignKey:FileMd5;references:Md5"`
	User        User        `gorm:"foreignKey:UserID"`
}

type FileMd5Info struct {
	Md5          string `gorm:"size:255;primary_key;not null;comment:文件 MD5"`
	Url          string `gorm:"unique;comment:文件 URL"`
	AbsolutePath string `gorm:"comment:文件绝对路径"`
}

func (*FileInfo) TableName() string {
	return FileInfoTable
}

func (*FileMd5Info) TableName() string {
	return FileInfoMd5Table
}
