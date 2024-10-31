package llm

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

// InputData 定义输入 JSON 的结构
type InputData struct {
	Model    string              `json:"model"`
	Messages []map[string]string `json:"messages"`
}

func NewChatService() *ChatService {
	return &ChatService{
		client: arkruntime.NewClientWithApiKey(
			os.Getenv("ARK_API_KEY"),
			arkruntime.WithBaseUrl("https://ark.cn-beijing.volces.com/api/v3"),
			arkruntime.WithRegion("cn-beijing"),
		),
	}
}

// ChatService 封装了 Volcengine 客户端
type ChatService struct {
	client *arkruntime.Client
}

// chatServiceInstance 保存单例实例
var chatServiceInstance *ChatService

// Init 使用环境变量中的 API key 初始化 ChatService 单例
func Init() error {
	apiKey := os.Getenv("ARK_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("缺少 API Key")
	}

	chatServiceInstance = &ChatService{
		client: arkruntime.NewClientWithApiKey(
			apiKey,
			arkruntime.WithBaseUrl("https://ark.cn-beijing.volces.com/api/v3"),
			arkruntime.WithRegion("cn-beijing"),
		),
	}

	return nil
}

// ProcessChatCompletion 处理聊天补全请求，返回 JSON 格式的 []byte 输出和可能的错误
func (cs *ChatService) ProcessChatCompletion(inputJSON []byte) ([]byte, error) {

	if chatServiceInstance == nil {
		return nil, fmt.Errorf("ChatService 未初始化，请先调用 Init()")
	}

	var inputData InputData
	if err := json.Unmarshal(inputJSON, &inputData); err != nil {
		return nil, fmt.Errorf("解析输入 JSON 错误: %v", err)
	}

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

	ctx := context.Background()
	resp, err := cs.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("请求错误: %v", err)
	}

	var nationality, confidenceLevel string
	if len(resp.Choices) > 0 && resp.Choices[0].Message.Content != nil {
		content := *resp.Choices[0].Message.Content.StringValue
		nationalityRegex := regexp.MustCompile(`\"possible_nation\"\s*:\s*\"([^\"]*)\"`) // 匹配英文国籍，例如 China、United States
		confidenceRegex := regexp.MustCompile(`\"confidence_level\"\s*:\s*(\d+)`)        // 匹配百分比置信度，例如 80%

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

	outputJSON, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("生成输出 JSON 错误: %v", err)
	}

	return outputJSON, nil
}
