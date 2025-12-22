package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OllamaEmbeddingService 本地Ollama向量化服务
type OllamaEmbeddingService struct {
	baseURL string
	model   string
	client  *http.Client
}

// NewOllamaEmbeddingService 创建Ollama向量化服务
func NewOllamaEmbeddingService() *OllamaEmbeddingService {
	return &OllamaEmbeddingService{
		baseURL: "http://localhost:11434",
		model:   "bge-m3",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// OllamaEmbeddingRequest Ollama API请求
type OllamaEmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// OllamaEmbeddingResponse Ollama API响应
type OllamaEmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

// GetEmbedding 获取文本的嵌入向量
func (s *OllamaEmbeddingService) GetEmbedding(text string) ([]float64, error) {
	if text == "" {
		return nil, fmt.Errorf("文本不能为空")
	}

	// 构建请求
	request := OllamaEmbeddingRequest{
		Model:  s.model,
		Prompt: text,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 调用Ollama API
	url := fmt.Sprintf("%s/api/embeddings", s.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("调用Ollama API失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		// 如果是404，可能是模型不存在
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("模型不存在，请执行: ollama pull bge-m3")
		}
		return nil, fmt.Errorf("Ollama返回错误 (%d): %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var result OllamaEmbeddingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return result.Embedding, nil
}

// TestConnection 测试Ollama连接
func (s *OllamaEmbeddingService) TestConnection() error {
	// 发送一个简单的测试请求
	_, err := s.GetEmbedding("test")
	if err != nil {
		return fmt.Errorf("Ollama连接测试失败: %v", err)
	}
	return nil
}
