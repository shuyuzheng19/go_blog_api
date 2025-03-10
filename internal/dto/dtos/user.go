package dtos

type UserLoginStatus struct {
	LastLogin int64  `gorm:"column:last_login"`
	LoginIp   string `gorm:"column:login_ip"`
	LoginCity string `gorm:"column:login_city"`
	ID        int    `gorm:"column:id"` // 这里可以省略，因为 GORM 默认会将 ID 作为主键处理
}
