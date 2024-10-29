package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/acd19ml/TalentRank/apps/user"
	"github.com/acd19ml/TalentRank/apps/user/impl"
)

func NewLLMJsonHandler() *Handler {
	svc := impl.NewUserServiceImpl()
	if svc == nil {
		panic("Failed to initialize UserServiceImpl")
	}
	return &Handler{svc: svc}
}

// Service 接口的实例
type Handler struct {
	svc user.Service
}

// 定义 Message 和 Request 结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// GetUserReposJSON 调用 CreateUserRepos 并返回 JSON 响应
func (h *Handler) GetUserReposJSON(ctx context.Context, username string) (string, error) {
	if username == "" {
		return "", fmt.Errorf("username cannot be empty")
	}

	// 调用 CreateUserRepos 获取 UserRepos 数据
	userRepos, err := h.svc.CreateUserRepos(ctx, username)
	if err != nil {
		return "", fmt.Errorf("failed to create user repos: %w", err)
	}

	// 将 UserRepos 编码为 JSON
	jsonData, err := json.Marshal(userRepos)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user repos to JSON: %w", err)
	}

	return string(jsonData), nil
}

// GetUserReposJSONWithRequest 构造请求数据并返回 JSON
func (h *Handler) GetUserReposJSONWithRequest(ctx context.Context, username string) (string, error) {
	// 调用 GetUserReposJSON 以获取用户仓库数据的 JSON 字符串
	userReposJSON, err := h.GetUserReposJSON(ctx, username)
	if err != nil {
		return "", err
	}

	// 构造消息内容
	request := Request{
		Model: "Model : ep-20241028212010-j7fg5",
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是豆包，是由字节跳动开发的 AI 人工智能助手",
			},
			{
				Role:    "user",
				Content: "Here is the information of GitHub user. Based on this information, you must infer their nation and provide a confidence level, and only respond with the nation and confidence level, the confidence level should be number in percent" + userReposJSON,
			},
		},
	}

	// 将请求数据编码为 JSON
	finalJSON, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request to JSON: %w", err)
	}

	return string(finalJSON), nil
}
