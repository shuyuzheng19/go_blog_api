package store

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

// saveToFile 将文件从源复制到目标路径
func SaveToFile(file *multipart.FileHeader, dstPath string) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("打开源文件失败: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("复制文件失败: %v", err)
	}

	return nil
}
