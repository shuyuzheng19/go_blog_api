package search

import (
	"blog/internal/utils"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type MeiliSearchClient struct {
	uri     string
	headers http.Header
}

// NewMeiliSearchClient 创建一个新的 MeiliSearch 客户端
func NewMeiliSearchClient(host string, apiKey string) *MeiliSearchClient {
	headers := http.Header{}
	headers.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	headers.Add("Content-Type", "application/json;charset=utf-8")

	return &MeiliSearchClient{
		uri:     host,
		headers: headers,
	}
}

// SendRequest 发送 HTTP 请求
func (c *MeiliSearchClient) SendRequest(method string, endpoint string, body string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.uri, endpoint)

	client := http.Client{}
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, err // 创建请求失败
	}

	request.Header = c.headers // 设置请求头

	response, err := client.Do(request) // 发送请求
	if err != nil {
		return nil, err // 请求失败
	}

	return response, nil
}

// CreateIndex 创建索引
func (c *MeiliSearchClient) CreateIndex(index string) error {
	var endpoint = "indexes"

	// 构建索引的负载
	var payload = map[string]interface{}{
		"uid":        index,
		"primaryKey": "id",
	}

	var jsonStr = utils.Serialize(payload) // 序列化为 JSON 字符串

	// 发送请求创建索引
	_, err := c.SendRequest(http.MethodPost, endpoint, jsonStr)

	return err
}

// DropIndex 删除索引
func (c *MeiliSearchClient) DropIndex(index string) error {
	var endpoint = "indexes/" + index

	// 发送请求删除索引
	_, err := c.SendRequest(http.MethodDelete, endpoint, "")

	return err
}

// DeleteAllDocument 删除索引中的所有文档
func (c *MeiliSearchClient) DeleteAllDocument(index string) error {
	var endpoint = fmt.Sprintf("indexes/%s/documents", index)

	// 发送请求删除所有文档
	_, err := c.SendRequest(http.MethodDelete, endpoint, "")

	return err
}

// SaveDocument 保存文档到索引
func (c *MeiliSearchClient) SaveDocument(index string, jsonDocument string) error {
	var endpoint = fmt.Sprintf("indexes/%s/documents", index)

	// 发送请求保存文档
	_, err := c.SendRequest(http.MethodPost, endpoint, jsonDocument)

	return err
}

// SearchDocument 在索引中搜索文档
func (c *MeiliSearchClient) SearchDocument(index string, req MeiliSearchRequest) MeiliSearchResponse {
	var endpoint = fmt.Sprintf("indexes/%s/search", index)

	// 发送搜索请求
	response, err := c.SendRequest(http.MethodPost, endpoint, utils.Serialize(req))
	if err != nil {
		return MeiliSearchResponse{} // 请求失败，返回空响应
	}
	defer response.Body.Close() // 确保在函数结束时关闭响应体

	if response.StatusCode == http.StatusOK {
		// 读取响应体
		result, _ := io.ReadAll(response.Body)
		return utils.Deserialize[MeiliSearchResponse](string(result)) // 反序列化响应
	}

	return MeiliSearchResponse{} // 返回空响应
}
