package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/acd19ml/TalentRank/apps/llm"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

func (s *LLMServer) ProcessChatCompletion(ctx context.Context, req *llm.ChatRequest) (*llm.ChatResponse, error) {
	// 使用 chatServiceInstance 来完成请求
	output, err := s.processCompletion([]byte(req.InputJson))
	if err != nil {
		return nil, fmt.Errorf("failed to process chat completion: %w", err)
	}

	return &llm.ChatResponse{OutputJson: string(output)}, nil
}

// processCompletion 为 JSON 处理的核心逻辑
func (s *LLMServer) processCompletion(inputJSON []byte) ([]byte, error) {
	var inputData InputData
	if err := json.Unmarshal(inputJSON, &inputData); err != nil {
		return nil, fmt.Errorf("解析输入 JSON 错误: %v", err)
	}

	var messages []*model.ChatCompletionMessage
	for _, msg := range inputData.Messages {
		roleEnum := map[string]string{
			"system":    model.ChatMessageRoleSystem,
			"user":      model.ChatMessageRoleUser,
			"assistant": model.ChatMessageRoleAssistant,
		}[msg["role"]]

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
	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("请求错误: %v", err)
	}

	var nationality, confidenceLevel string
	if len(resp.Choices) > 0 && resp.Choices[0].Message.Content != nil {
		content := *resp.Choices[0].Message.Content.StringValue
		nationalityRegex := regexp.MustCompile(`\"possible_nation\"\s*:\s*\"([^\"]*)\"`)
		confidenceRegex := regexp.MustCompile(`\"confidence_level\"\s*:\s*(\d+)`)

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
