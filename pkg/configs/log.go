package configs

import (
	"blog/pkg/helper"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerConfig 日志配置结构体
type LoggerConfig struct {
	Dev         bool   `yaml:"dev" json:"dev"`                 // 是否开发模式
	Encoding    string `yaml:"encoding" json:"encoding"`       // 日志编码格式
	OutputPaths string `yaml:"outputPaths" json:"outputPaths"` // 输出路径
	ErrorPaths  string `yaml:"errorPaths" json:"errorPaths"`   // 错误输出路径
	Level       string `yaml:"level" json:"level"`             // 日志级别
	LoggerDir   string `yaml:"loggerDir" json:"loggerDir"`     // 日志文件目录
	DefaultName string `yaml:"defaultName" json:"defaultName"` // 默认日志文件名
}

// 日志级别映射
var levelMap = map[string]zapcore.Level{
	"info":  zapcore.InfoLevel,
	"debug": zapcore.DebugLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"panic": zapcore.PanicLevel,
}

var LOGGER *zap.Logger // 全局日志实例

// LoadLogger 加载 zap 日志配置
func LoadLogger(loggerConfig LoggerConfig) *zap.Logger {
	// 检查默认日志文件名是否为空
	if !loggerConfig.Dev && loggerConfig.DefaultName == "" {
		helper.CheckError(errors.New("log defaultName empty"), "日志默认文件名不能为空")
	}

	// 创建日志目录
	if err := os.MkdirAll(loggerConfig.LoggerDir, os.ModePerm); err != nil {
		helper.CheckError(err, "创建日志目录失败")
	}

	// 创建 zap 日志配置
	config := zap.NewProductionConfig()
	config.Encoding = loggerConfig.Encoding
	config.Development = loggerConfig.Dev

	// 设置日志级别
	level, found := levelMap[loggerConfig.Level]
	if !found {
		level = zap.InfoLevel // 默认日志级别为 Info
	}
	config.Level = zap.NewAtomicLevelAt(level)

	// 设置时间格式
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	// 设置输出路径
	paths := []string{}
	if loggerConfig.OutputPaths != "" {
		paths = strings.Split(loggerConfig.OutputPaths, ",")
	}
	// 添加默认日志文件路径
	defaultLogPath := filepath.Join(loggerConfig.LoggerDir, loggerConfig.DefaultName)
	paths = append(paths, defaultLogPath)
	config.OutputPaths = paths

	// 自动生成错误日志文件路径
	if loggerConfig.ErrorPaths == "" {
		errorLogPath := filepath.Join(loggerConfig.LoggerDir, "error.log")
		config.ErrorOutputPaths = []string{errorLogPath}
	} else {
		config.ErrorOutputPaths = strings.Split(loggerConfig.ErrorPaths, ",")
	}

	// 构建 logger 实例
	logger, err := config.Build()
	if err != nil {
		helper.CheckError(err, "加载日志配置失败")
	}

	return logger
}
