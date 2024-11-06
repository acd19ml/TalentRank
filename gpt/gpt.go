package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// 定义 OpenAI 请求结构体
type OpenAIRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 定义 OpenAI 响应结构体
type OpenAIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// 调用 OpenAI API
func callOpenAI(prompt string) (string, error) {
	apiKey := os.Getenv("GPT_API") // 确保将API密钥存入环境变量
	if apiKey == "" {
		log.Fatal("API key not found")
	}
	url := "https://api.openai.com/v1/chat/completions"

	// 构建请求数据
	requestData := OpenAIRequest{
		Model: "gpt-4o", // 使用 GPT-4 模型
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant to guess the nation and confidence_level for a GitHub user by multiple data from this user."},
			{Role: "user", Content: prompt},
		},
		MaxTokens: 100,
	}

	// 转换为 JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

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

	// 返回生成的文本
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("no response from OpenAI")
}

// 处理 HTTP 请求
func openAIHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Prompt string `json:"prompt"`
	}

	// 解码请求体中的 JSON
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 调用 OpenAI API
	response, err := callOpenAI(request.Prompt)
	if err != nil {
		http.Error(w, "Error calling OpenAI: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回 JSON 响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // 默认端口为 8080
	}
	http.HandleFunc("/analyze", openAIHandler)
	log.Println("Server is running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
