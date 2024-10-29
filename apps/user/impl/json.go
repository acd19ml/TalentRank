package impl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/acd19ml/TalentRank/apps/user"
)

// 定义 Message 和 Request 结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// GetUserReposJSONWithRequest 构造请求数据并返回 JSON
func GetUserReposJSONWithRequest(ctx context.Context, user *user.User) ([]byte, error) {
	// 调用 GetUserReposJSON 以获取用户仓库数据的 JSON 字符串
	userJSON, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user repos to JSON: %w", err)
	}

	// 构造消息内容
	request := Request{
		Model: "ep-20241029222143-j4r4t",
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是豆包，是由字节跳动开发的 AI 人工智能助手",
			},
			{
				Role:    "user",
				Content: "Here is the information of GitHub user. Based on this information, you must infer their nation and provide a confidence level, and only respond with the nation and confidence level, the confidence level should be number in percent" + string(userJSON),
			},
		},
	}

	// 将请求数据编码为 JSON
	finalJSON, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request to JSON: %w", err)
	}

	return finalJSON, nil
}
