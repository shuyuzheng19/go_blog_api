package service

import (
	"blog/internal/dto/requests"
	"blog/internal/dto/response"
	"blog/internal/job"
	"blog/pkg/configs"
	"blog/pkg/logger"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
)

type SystemService struct {
}

// GetSystemFile 获取系统文件列表
func (fs *SystemService) GetSystemFile(req requests.SystemFileRequest) []response.SystemFileResponse {
	if req.Path == "" {
		req.Path = configs.CONFIG.Upload.Path
	}

	var result []response.SystemFileResponse

	// 检查路径是否存在
	if _, err := os.Stat(req.Path); os.IsNotExist(err) {
		return result
	}

	// 读取目录文件
	files, err := os.ReadDir(req.Path) // 使用 os.ReadDir 替代 ioutil.ReadDir
	if err != nil {
		return result
	}

	for _, info := range files {
		if info.IsDir() {
			continue
		}

		name := info.Name()
		// 使用 strings.EqualFold 进行不区分大小写的比较
		if req.Keyword != "" && !strings.Contains(strings.ToLower(name), strings.ToLower(req.Keyword)) {
			continue
		}

		// 获取文件的详细信息
		filePath := filepath.Join(req.Path, name)
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			continue // 如果获取文件信息失败，跳过该文件
		}

		result = append(result, response.SystemFileResponse{
			Name:       name,
			Path:       filePath,
			Ext:        filepath.Ext(name),
			Size:       fileInfo.Size(),
			CreateTime: fileInfo.ModTime().Unix(), // 创建时间
			UpdateTime: fileInfo.ModTime().Unix(), // 修改时间
		})
	}

	return result
}

// DeleteSystemFile 删除系统文件
func (fs *SystemService) DeleteSystemFile(paths []string) int64 {
	var count int64
	for _, path := range paths {
		if err := os.Remove(path); err == nil {
			count++
		} else {
			// 记录错误信息（可选）
			// log.Printf("Failed to delete file: %s, error: %v", path, err)
		}
	}
	return count
}

// ClearFileContent 清空文件内容
func (fs *SystemService) ClearFileContent(path string) error {
	return os.Truncate(path, 0)
}

func (fs *SystemService) autoCreateFileLog() {
	var logConfig = configs.CONFIG.Logger

	var logPath = filepath.Join(logConfig.LoggerDir, logConfig.DefaultName)

	var fileName = time.Now().Add(-time.Minute).Format("2006-01-02") + ".log"

	var file, err = os.OpenFile(filepath.Join(logConfig.LoggerDir, fileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)

	defer file.Close()

	if err != nil {
		logger.Info("创建日志文件失败", zap.String("error", err.Error()))
		return
	}

	fs.ClearFileContent(logPath)

	var buff, _ = os.ReadFile(logPath)

	file.WriteString(string(buff))

	logger.Info("已重新创建日志", zap.String("file_name", fileName))
}

// NewSystemService 创建新的 SystemService 实例
func NewSystemService() *SystemService {
	var service = &SystemService{}

	if configs.CONFIG.Server.Cron {
		job.AddJob(job.Job{
			Hour:        0,
			Eq:          true,
			Description: "自动创建日志",
			Job:         service.autoCreateFileLog,
		})
	}

	return service
}
