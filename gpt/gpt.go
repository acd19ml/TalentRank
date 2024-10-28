package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type OpenAIRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func callOpenAI(prompt string) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	apiKey := os.Getenv("GPT_API") // 确保将API密钥存入环境变量
	fmt.Println("Loaded API Key:", apiKey)
	url := "https://api.openai.com/v1/chat/completions"

	// 构造请求数据
	requestData := OpenAIRequest{
		Model: "gpt-3.5-turbo", // 使用GPT-3.5-turbo模型
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
		MaxTokens: 100, // 设置返回内容的最大字符数
	}

	// 将请求数据转换为JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// client := &http.Client{
	// 	Timeout: 30 * time.Second, // 将超时时间设置为30秒或更长
	// }

	// 发送请求并获取响应
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析响应
	var response OpenAIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	// 返回模型生成的文本
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from OpenAI")
}

func main() {
	response, err := callOpenAI("你好")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Response:", response)
	}
}
