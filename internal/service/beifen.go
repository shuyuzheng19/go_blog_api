package service

import (
	"blog/pkg/common"
	"blog/pkg/configs"
	"blog/pkg/logger"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// 自动删除文件的方法
func scheduleFileDeletion(filePath string, delay time.Duration) {
	time.AfterFunc(delay, func() {
		err := os.Remove(filePath)
		if err != nil {
			logger.Info("文件删除失败", zap.String("path", filePath), zap.String("err", err.Error()))
		} else {
			logger.Info("文件已删除", zap.String("path", filePath))
		}
	})
}

func CreateTarGz(sourceDir, outputFile string) error {
	cmd := exec.Command("tar", "-czf", outputFile, sourceDir)
	return cmd.Run()
}

// 打包目录到父级目录，并使用UUID命名的tar包
func tarDirectoryWithUUID(uuid, sourceDir string) (string, error) {
	// 获取父级目录
	parentDir := filepath.Dir(sourceDir)

	tarFile := filepath.Join(parentDir, uuid+".tar.gz")

	// 执行tar命令，打包目录到tar.gz文件
	cmd := exec.Command("tar", "-czf", tarFile, "-C", parentDir, filepath.Base(sourceDir))
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("打包失败: %v", err)
	}

	return tarFile, nil
}

// 创建打包并存入Redis
func CreateTarGzAndStoreRedis(path string, min int) (string, error) {

	// 生成一个UUID作为压缩包的文件名
	uniqueID := uuid.New().String()

	var filePath, err = tarDirectoryWithUUID(uniqueID, path)

	if err != nil {
		log.Info("压缩tar包失败", zap.String("err", err.Error()))
		return "", err
	}

	if err := configs.REDIS.SetNX(common.TarKey+uniqueID, filePath, time.Minute*time.Duration(min)).Err(); err != nil {
		return "", err
	}

	logger.Info("存入打包路径到redis", zap.String("path", filePath))
	// 调用定时删除
	logger.Info("调用定时删除", zap.Int("min", min))

	go scheduleFileDeletion(filePath, time.Duration(min)*time.Minute)

	return uniqueID, nil

}
