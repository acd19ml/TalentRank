package impl

import (
	"context"
	"encoding/json"
	"errors"
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
	if user == nil {
		return nil, errors.New("userins is nil in GetUserReposJSONWithRequest")
	}
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
				Content: "You are Doubao, an AI assistant developed by ByteDance.",
			},
			{
				Role:    "user",
				Content: "Based on the provided JSON information of a GitHub user (id, username, name, company, blog, location, email, bio, followers, organizations, readme, commits, and score), if the location field is empty, analyze and infer the user’s possible nation (location) based on the other available fields. Only respond with a JSON object in the following format: \n\n```\n{\n  \"possible_nation\": \"<country or 'N/A'>\",\n  \"confidence_level\": <percentage as a number>\n}\n```\n\nIf the information is insufficient to determine the user's nation, set \"possible_nation\" to \"N/A\" and \"confidence_level\" to 0." + string(userJSON),
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
