package configs

import (
	"blog/pkg/helper"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var CONFIG *GlobalConfig

// GlobalConfig 全局配置
type GlobalConfig struct {
	//环境
	Env string `yaml:"env" json:"-"`
	// IP数据库路径
	IpDbPath string `yaml:"ipDbPath" json:"ipDbPath"`
	// server配置
	Server ServerConfig `yaml:"server" json:"server"`
	//db配置
	Db DbConfig `yaml:"db" json:"db"`
	//日志配置
	Logger LoggerConfig `yaml:"logger" json:"logger"`
	//redis配置
	Redis RedisConfig `yaml:"redis" json:"redis"`
	//SMTP配置
	Mail EmailConfig `yaml:"email" json:"email"`
	//我的邮箱
	MyEmail string `yaml:"myEmail" json:"myEmail"`
	//上传文件配置
	Upload UploadConfig `yaml:"upload" json:"upload"`
	//搜索配置
	Search      MeiliSearchConfig `yaml:"meilisearch" json:"meilisearch"`
	DataBaseKey string            `yaml:"databaseKey" json:"-"`
}

// LoadGlobalConfig 加载全局配置
func LoadGlobalConfig(path string) {
	var file, err = os.ReadFile(path)

	helper.CheckError(err, "读取配置文件失败")

	yaml.Unmarshal(file, &CONFIG)

	fmt.Println("config=", CONFIG)
}
