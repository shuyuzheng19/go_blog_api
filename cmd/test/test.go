package main

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
)

func tarDirectoryWithUUID(sourceDir string) (string, error) {
	// 获取父级目录
	parentDir := filepath.Dir(sourceDir)

	// 生成一个UUID作为压缩包的文件名
	uniqueID := uuid.New().String()
	tarFile := filepath.Join(parentDir, uniqueID+".tar.gz")

	// 执行tar命令，打包目录到tar.gz文件
	cmd := exec.Command("tar", "-czf", tarFile, "-C", parentDir, filepath.Base(sourceDir))
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("打包失败: %v", err)
	}

	return tarFile, nil
}

func main() {
	tarDirectoryWithUUID("/zsy_test")
}
