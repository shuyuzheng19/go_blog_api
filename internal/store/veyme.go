package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type VeyMeResponse struct {
	Image string `json:"image"`
	Url   string `json:"url"`
}

func UploadImageToVeyme(file *multipart.FileHeader, token string) (*VeyMeResponse, error) {
	// 打开文件
	var f, err = file.Open()

	if err != nil {
		return nil, err
	}

	defer f.Close()

	// 创建一个缓冲区来存储表单数据
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 创建文件表单字段
	part, err := writer.CreateFormFile("file", "image.jpg")
	if err != nil {
		fmt.Println("创建表单文件字段时发生错误:", err)
		return nil, err
	}

	// 将文件内容复制到表单字段
	_, err = io.Copy(part, f)
	if err != nil {
		fmt.Println("复制文件内容时发生错误:", err)
		return nil, err
	}

	// 添加用户密钥字段
	err = writer.WriteField("userkey", token)
	if err != nil {
		fmt.Println("写入用户密钥字段时发生错误:", err)
		return nil, err
	}

	// 关闭 writer，完成表单数据的构建
	err = writer.Close()
	if err != nil {
		fmt.Println("关闭 writer 时发生错误:", err)
		return nil, err
	}

	// 创建一个 POST 请求
	req, err := http.NewRequest("POST", "https://vgy.me/upload", &requestBody)
	if err != nil {
		fmt.Println("创建请求时发生错误:", err)
		return nil, err
	}

	// 设置 Content-Type 为 multipart 表单的内容类型
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求时发生错误:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应时发生错误:", err)
		return nil, err
	}

	var image VeyMeResponse

	json.Unmarshal(body, &image)

	// 输出响应内容
	fmt.Println(string(body))

	return &image, nil
}
