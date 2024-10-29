package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

type InputData struct {
	Model    string              `json:"model"`
	Messages []map[string]string `json:"messag:qes"`
}

func main() {
	inputFile, err := os.Open("./doubao/input.json")
	if err != nil {
		fmt.Printf("无法打开输入文件: %v\n", err)
		return
	}
	defer inputFile.Close()

	var inputData InputData
	if err := json.NewDecoder(inputFile).Decode(&inputData); err != nil {
		fmt.Printf("解析输入文件错误: %v\n", err)
		return
	}

	client := arkruntime.NewClientWithApiKey(
		os.Getenv("ARK_API_KEY"),
		arkruntime.WithBaseUrl("https://ark.cn-beijing.volces.com/api/v3"),
		arkruntime.WithRegion("cn-beijing"),
	)

	ctx := context.Background()

	var messages []*model.ChatCompletionMessage
	for _, msg := range inputData.Messages {
		var roleEnum string
		switch msg["role"] {
		case "system":
			roleEnum = model.ChatMessageRoleSystem
		case "user":
			roleEnum = model.ChatMessageRoleUser
		case "assistant":
			roleEnum = model.ChatMessageRoleAssistant
		}

		messages = append(messages, &model.ChatCompletionMessage{
			Role: roleEnum,
			Content: &model.ChatCompletionMessageContent{
				StringValue: volcengine.String(msg["content"]),
			},
		})
	}

	req := model.ChatCompletionRequest{
		Model:    inputData.Model,
		Messages: messages,
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("请求错误: %v\n", err)
		return
	}

	var nationality, confidenceLevel string
	if len(resp.Choices) > 0 && resp.Choices[0].Message.Content != nil {
		content := *resp.Choices[0].Message.Content.StringValue
		print(content)
		// 使用正则表达式提取中文字段“国籍”和“信心等级”
		// nationalityRegex := regexp.MustCompile(`"possible nation":\s*"([^"]+)"`)
		// confidenceRegex := regexp.MustCompile(`"confidence level":\s*([0-9.]+)`)
		// nationalityRegex := regexp.MustCompile(`"国籍":\s*"([^"]+)"`)
		// confidenceRegex := regexp.MustCompile(`"置信度":\s*([0-9.]+)`)
		// nationalityRegex := regexp.MustCompile(`([\p{Han}]+)`) // 匹配中国、美国等中文国籍
		// confidenceRegex := regexp.MustCompile(`([0-9]+)%`)     // 匹配百分比置信度
		nationalityRegex := regexp.MustCompile(`([A-Za-z]+)`) // 匹配英文国籍，例如 China、United States
		confidenceRegex := regexp.MustCompile(`([0-9]+)%`)    // 匹配百分比置信度，例如 80%

		nationalityMatch := nationalityRegex.FindStringSubmatch(content)
		confidenceMatch := confidenceRegex.FindStringSubmatch(content)

		if len(nationalityMatch) > 1 {
			nationality = nationalityMatch[1]
		}
		if len(confidenceMatch) > 1 {
			confidenceLevel = confidenceMatch[1]
		}
	}

	output := map[string]interface{}{
		"possible_nation":  nationality,
		"confidence_level": confidenceLevel,
	}

	outputFile, err := os.Create("output.json")
	if err != nil {
		fmt.Printf("无法创建输出文件: %v\n", err)
		return
	}
	defer outputFile.Close()

	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(output); err != nil {
		fmt.Printf("写入输出文件错误: %v\n", err)
		return
	}

	fmt.Println("简化响应已保存到 output.json")
}
