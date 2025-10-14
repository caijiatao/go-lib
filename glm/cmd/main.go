package main

import (
	"fmt"
	"golib/glm"
	"log"
)

func main() {
	// TODO API key
	apiKey := ""
	// 创建客户端
	client := glm.NewGLMClient(apiKey)

	// 示例1: 非流式调用
	fmt.Println("=== 非流式调用 ===")
	messages := []glm.Message{
		{
			Role:    "system",
			Content: "你是一个有用的AI助手。",
		},
		{
			Role:    "user",
			Content: "你好，请介绍一下自己。",
		},
	}

	result, err := client.SimpleChat("glm-4", messages, 0.6)
	if err != nil {
		log.Fatalf("非流式调用失败: %v", err)
	}
	fmt.Printf("回复: %s\n\n", result)

	// 示例2: 流式调用
	fmt.Println("=== 流式调用 ===")
	streamReq := &glm.ChatRequest{
		Model:       "glm-4",
		Messages:    messages,
		Temperature: 0.6,
		Stream:      true,
	}

	fmt.Print("流式回复: ")
	err = client.ChatStream(streamReq, func(chunk *glm.StreamResponseChunk) error {
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("流式调用失败: %v", err)
	}
	fmt.Println("\n")

	// 示例3: 完整功能调用
	fmt.Println("=== 完整功能调用 ===")
	fullReq := &glm.ChatRequest{
		Model: "glm-4",
		Messages: []glm.Message{
			{
				Role:    "system",
				Content: "你是一个编程专家。",
			},
			{
				Role:    "user",
				Content: "用Go写一个快速排序算法。",
			},
		},
		Temperature: 0.7,
		MaxTokens:   1000,
	}

	fullResult, err := client.Chat(fullReq)
	if err != nil {
		log.Fatalf("完整功能调用失败: %v", err)
	}
	fmt.Printf("回复: %s\n", fullResult.Content)
	fmt.Printf("Token使用: 提示=%d, 完成=%d, 总计=%d\n",
		fullResult.Usage.PromptTokens,
		fullResult.Usage.CompletionTokens,
		fullResult.Usage.TotalTokens)
}
