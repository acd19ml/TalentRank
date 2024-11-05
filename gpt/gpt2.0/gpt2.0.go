package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
)

func main() {

	// 获取 API 密钥
	apiKey := os.Getenv("GPT_API")
	if apiKey == "" {
		log.Fatal("请设置GPT_API环境变量")
	}

	// 自定义 HTTP 客户端，添加 Authorization 头部
	httpClient := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &transportWithAuth{
			apiKey: apiKey,
			rt:     http.DefaultTransport,
		},
	}

	// 使用自定义的 HTTP 客户端和 BaseURL
	client := openai.NewClientWithConfig(openai.ClientConfig{
		BaseURL:    "https://api.openai.com/v1",
		HTTPClient: httpClient,
	})

	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model: "gpt-4",
		Messages: []openai.ChatCompletionMessage{
			{Role: "user", Content: "写一篇《时间管理》"},
		},
		MaxTokens:   2048,
		Temperature: 0.5,
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		log.Fatalf("调用 OpenAI API 出错: %v", err)
	}

	fmt.Println(resp.Choices[0].Message.Content)
}

// 自定义 Transport 结构，用于设置 Authorization 头部
type transportWithAuth struct {
	apiKey string
	rt     http.RoundTripper
}

func (t *transportWithAuth) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+t.apiKey)
	return t.rt.RoundTrip(req)
}
