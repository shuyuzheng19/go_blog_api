package utils

import (
	"crypto/md5"
	"encoding/hex"
	"mime/multipart"

	"io"
)

// CalculateMD5 计算文件的MD5值
func CalculateMD5(fileHeader *multipart.FileHeader) string {
	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return ""
	}
	defer file.Close() // 确保在函数结束时关闭文件

	hash := md5.New()
	buff := make([]byte, 4096) // 使用4KB的缓冲区

	// 逐块读取文件并更新MD5哈希
	for {
		n, err := file.Read(buff)
		if err != nil {
			if err == io.EOF {
				break // 到达文件末尾
			}
			return "" // 读取错误
		}
		hash.Write(buff[:n]) // 只写入实际读取的字节
	}

	return hex.EncodeToString(hash.Sum(nil)) // 返回MD5值
}
