package configs

import (
	"blog/internal/models"
	"blog/internal/utils"
	"blog/pkg/common"
	"blog/pkg/helper"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// DbConfig 关系型数据库配置
type DbConfig struct {
	Log        bool   `yaml:"log" json:"log"`               // 是否开启日志
	MaxIdle    int    `yaml:"maxIdle" json:"maxIdle"`       // 空闲连接数
	MaxSize    int    `yaml:"maxSize" json:"maxSize"`       // 最大连接数
	Timezone   string `yaml:"timezone" json:"timezone"`     // 数据库时区
	Database   string `yaml:"database" json:"database"`     // 数据库类型，如 mysql、postgresql 等
	Host       string `yaml:"host" json:"host"`             // 数据库主机
	Port       int    `yaml:"port" json:"port"`             // 数据库端口
	Username   string `yaml:"username" json:"username"`     // 数据库用户名
	Password   string `yaml:"password" json:"password"`     // 数据库密码
	Dbname     string `yaml:"dbname" json:"dbname"`         // 数据库名称
	AutoCreate bool   `yaml:"autoCreate" json:"autoCreate"` // 是否自动创建表
}

// getDataBaseDSN 根据 DbConfig 返回数据库连接字符串
func getDataBaseDSN(config DbConfig) string {
	var database = strings.ToLower(config.Database) // 将数据库类型转换为小写
	switch database {
	case "mysql":
		// MySQL 数据库连接字符串
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=%s",
			config.Username, config.Password, config.Host, config.Port, config.Dbname, config.Timezone)
	case "postgresql":
		// PostgreSQL 数据库连接字符串
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable Timezone=%s",
			config.Host, config.Port, config.Username, config.Dbname, config.Password, config.Timezone)
	case "sqlite":
		// SQLite 数据库连接字符串
		return config.Dbname // 假设 Dbname 是 SQLite 文件的路径
	case "sqlserver":
		// SQL Server 数据库连接字符串
		return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
			config.Username, config.Password, config.Host, config.Port, config.Dbname)
	case "oracle":
		// Oracle 数据库连接字符串
		return fmt.Sprintf("%s/%s@%s:%d/%s",
			config.Username, config.Password, config.Host, config.Port, config.Dbname)
	case "cockroachdb":
		// CockroachDB 数据库连接字符串
		return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
			config.Username, config.Password, config.Host, config.Port, config.Dbname)
	case "clickhouse":
		// ClickHouse 数据库连接字符串
		return fmt.Sprintf("tcp://%s:%d?username=%s&password=%s&database=%s",
			config.Host, config.Port, config.Username, config.Password, config.Dbname)
	case "bigquery":
		// BigQuery 数据库连接字符串
		return fmt.Sprintf("bigquery://%s:%s@projectID:%s/datasetID",
			config.Username, config.Password, config.Dbname)
	default:
		// 如果数据库类型未知，抛出错误
		panic(fmt.Sprintf("未知的数据库 %s", config.Database))
	}
}

// LoadDBConfig 加载数据库配置
func LoadDBConfig(dbConfig DbConfig) {
	// 获取数据库连接字符串
	var dsn = getDataBaseDSN(dbConfig)

	// 创建 GORM 配置
	var gormConfig = &gorm.Config{}

	// 根据配置决定是否开启日志
	if dbConfig.Log {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	// 打开数据库连接
	var db, err = gorm.Open(postgres.Open(dsn), gormConfig)
	// 判断是否连接成功
	helper.CheckError(err, "数据库加载失败")

	// 获取数据库连接池
	var connection, _ = db.DB()
	// 设置最大空闲连接数
	connection.SetMaxIdleConns(dbConfig.MaxIdle)
	// 设置最大连接数
	connection.SetMaxOpenConns(dbConfig.MaxSize)

	DB = db // 将数据库实例赋值给全局变量

	if dbConfig.AutoCreate {
		DB.AutoMigrate(&models.User{},
			&models.Role{},
			&models.FileInfo{},
			&models.Blog{},
			&models.EyeView{},
			&models.SystemLogInfo{},
			&models.EditBlog{},
		)

		var roles = []models.Role{
			{
				ID:          uint(common.UserRoleId),
				Name:        "USER",
				Description: "普通用户",
			},
			{
				ID:          uint(common.AdminRoleId),
				Name:        "ADMIN",
				Description: "管理员",
			},
			{
				ID:          uint(common.SuperAdminRoleId),
				Name:        "SUPER_ADMIN",
				Description: "超级管理员",
			},
		}

		DB.Model(&models.Role{}).Save(&roles)

		var user = models.User{
			Username: "2528959216",
			Password: utils.EncryptPassword("shuyu2001"),
			Email:    "shuyuzheng19@gmail.com",
			NickName: "书宇",
			Avatar:   "https://www.shuyuz.com/avatar.jpg",
			RoleID:   3,
			Status:   true,
		}

		DB.Model(&models.User{}).Save(&user)
	}
}

func DeleteData(tabeName string, uid *int, ids []int64) error {
	var db = DB.Table(tabeName).Where("id in ?", ids)

	if uid != nil {
		db.Where("user_id = ?", *uid)
	}

	return db.Update("deleted_at", time.Now()).Error
}

func UnDeleteData(tabeName string, uid *int, ids []int64) error {
	var db = DB.Table(tabeName).Where("id in ?", ids)

	if uid != nil {
		db.Where("user_id = ?", *uid)
	}

	return db.Update("deleted_at", nil).Error
}
