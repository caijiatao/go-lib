package glm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Message 表示对话消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 表示聊天请求
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// ChatResponse 表示非流式响应
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// StreamResponseChunk 表示流式响应的数据块
type StreamResponseChunk struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Delta struct {
			Role    string `json:"role,omitempty"`
			Content string `json:"content,omitempty"`
		} `json:"delta"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason,omitempty"`
	} `json:"choices"`
}

// ChatCompletion 表示完整的聊天完成结果
type ChatCompletion struct {
	Content string
	Usage   struct {
		PromptTokens     int
		CompletionTokens int
		TotalTokens      int
	}
}

// GLMClient 智谱AI客户端
type GLMClient struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// NewGLMClient 创建新的客户端实例
func NewGLMClient(apiKey string) *GLMClient {
	return &GLMClient{
		APIKey:     apiKey,
		BaseURL:    "https://open.bigmodel.cn/api/paas/v4",
		HTTPClient: &http.Client{},
	}
}

// SetBaseURL 设置自定义的基础URL
func (c *GLMClient) SetBaseURL(baseURL string) {
	c.BaseURL = baseURL
}

// Chat 非流式聊天
func (c *GLMClient) Chat(req *ChatRequest) (*ChatCompletion, error) {
	// 设置默认值
	if req.Model == "" {
		req.Model = "glm-4"
	}
	if req.Temperature == 0 {
		req.Temperature = 0.6
	}

	req.Stream = false

	url := fmt.Sprintf("%s/chat/completions", c.BaseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API返回错误: %s, 响应: %s", resp.Status, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("API返回空的选择")
	}

	result := &ChatCompletion{
		Content: chatResp.Choices[0].Message.Content,
		Usage: struct {
			PromptTokens     int
			CompletionTokens int
			TotalTokens      int
		}{
			PromptTokens:     chatResp.Usage.PromptTokens,
			CompletionTokens: chatResp.Usage.CompletionTokens,
			TotalTokens:      chatResp.Usage.TotalTokens,
		},
	}

	return result, nil
}

// ChatStream 流式聊天
func (c *GLMClient) ChatStream(req *ChatRequest, onChunk func(*StreamResponseChunk) error) error {
	// 设置默认值
	if req.Model == "" {
		req.Model = "glm-4"
	}
	if req.Temperature == 0 {
		req.Temperature = 0.6
	}

	req.Stream = true

	url := fmt.Sprintf("%s/chat/completions", c.BaseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %v", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API返回错误: %s, 响应: %s", resp.Status, string(body))
	}

	// 处理流式响应
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			// 流结束标记
			if data == "[DONE]" {
				break
			}

			var chunk StreamResponseChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				return fmt.Errorf("解析流数据失败: %v, 数据: %s", err, data)
			}

			if err := onChunk(&chunk); err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取流数据失败: %v", err)
	}

	return nil
}

// SimpleChat 简化版的聊天方法
func (c *GLMClient) SimpleChat(model string, messages []Message, temperature float64) (string, error) {
	req := &ChatRequest{
		Model:       model,
		Messages:    messages,
		Temperature: temperature,
	}

	result, err := c.Chat(req)
	if err != nil {
		return "", err
	}

	return result.Content, nil
}
