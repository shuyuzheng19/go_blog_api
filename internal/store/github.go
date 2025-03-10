package store

import (
	"blog/pkg/configs"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type GitHubContent struct {
	Message string `json:"message"`
	Content string `json:"content"`
	Branch  string `json:"branch,omitempty"`
}

func UploadImageToGitHub(file *multipart.FileHeader, path string, github *configs.GithubUploadConfig) (string, error) {

	var f, err = file.Open()

	if err != nil {
		return "", err
	}

	defer f.Close()

	// 读取图片文件
	imageData, err := io.ReadAll(f)

	if err != nil {
		return "", fmt.Errorf("failed to read image file: %v", err)
	}

	// 将图片内容进行 Base64 编码
	encodedContent := base64.StdEncoding.EncodeToString(imageData)

	// 构建请求体
	content := GitHubContent{
		Content: encodedContent,
		Branch:  "main", // 可选，默认是 "main" 或 "master"
	}

	// 将请求体编码为 JSON
	requestBody, err := json.Marshal(content)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// GitHub API 上传文件的 URL
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", github.User, github.Repo, path)

	// 创建 HTTP 请求
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", "token "+github.Token)
	req.Header.Set("Content-Type", "application/json")

	// 发送 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp == nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to upload image, status code: %d", resp.StatusCode)
	}

	var buff, _ = io.ReadAll(resp.Body)
	var r response
	if err := json.Unmarshal(buff, &r); err != nil {
		return "", err
	}
	// var urll = "https://site.yuflow.us.kg/" + r.Content.Path + "?ref=main&sha=" + r.Content.Sha
	// var urll = fmt.Sprintf("https://fastly.jsdelivr.net/gh/%s/%s@main/%s?sha=%s", owner, repo, r.Content.Path, r.Content.Sha)
	var urll = fmt.Sprintf(github.Proxy, github.User, github.Repo, r.Content.Path, r.Content.Sha)
	return urll, nil
}

type response struct {
	Content struct {
		Path        string `json:"path"`
		DownloadUrl string `json:"download_url"`
		Url         string `json:"url"`
		Sha         string `json:"sha"`
	} `json:"content"`
}
