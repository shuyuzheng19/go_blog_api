package logger

import (
	"blog/pkg/configs"

	"go.uber.org/zap"
)

// Logger 封装的日志实例
var Logger *zap.Logger

// InitLogger 初始化日志
func InitLogger(config configs.LoggerConfig) {
	Logger = configs.LoadLogger(config) // 调用之前的 LoadLogger 函数

	defer Logger.Sync()
}

// Info 记录信息级别日志
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Debug 记录调试级别日志
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Warn 记录警告级别日志
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Error 记录错误级别日志
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Panic 记录恐慌级别日志
func Panic(msg string, fields ...zap.Field) {
	Logger.Panic(msg, fields...)
}

// 明显日志
func InfoView(msg string) {
	Logger.Info("============================================= " + msg + " =============================================")
}

// Fatal 记录致命级别日志
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}
