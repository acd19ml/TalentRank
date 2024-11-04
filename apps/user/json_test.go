package user_test

// import (
// 	"context"
// 	"testing"

// 	"github.com/acd19ml/TalentRank/apps/llm/impl"
// 	"github.com/acd19ml/TalentRank/apps/user"
// )

// // TestProcessChatCompletionAndUnmarshal 测试 ProcessChatCompletion 和 UnmarshalToUserResponceByLLM 方法
// func TestProcessChatCompletionAndUnmarshal(t *testing.T) {
// 	// 初始化 ChatService 单例
// 	err := impl.Init()
// 	if err != nil {
// 		t.Fatalf("Failed to initialize ChatService: %v", err)
// 	}

// 	// 创建测试用的 user 数据
// 	ctx := context.Background()
// 	userins := &user.User{
// 		Username:      "testuser",
// 		Name:          "Test User",
// 		Location:      "", // 空 Location，触发 nation 和 confidence level 的推断
// 		Email:         "test@example.com",
// 		Bio:           "A test user",
// 		Followers:     42,
// 		Organizations: []string{"Org1", "Org2"},
// 	}

// 	// 生成请求 JSON
// 	finalJSON, err := user.GetUserReposJSONWithRequest(ctx, userins)
// 	if err != nil {
// 		t.Fatalf("GetUserReposJSONWithRequest failed: %v", err)
// 	}
// 	var svc user.LLMService
// 	svc = impl.NewChatService()
// 	// 调用 ProcessChatCompletion 并获取响应
// 	responseJSON, err := svc.ProcessChatCompletion(finalJSON)
// 	if err != nil {
// 		t.Fatalf("ProcessChatCompletion failed: %v", err)
// 	}

// 	// 使用 UnmarshalToUserResponceByLLM 解析响应
// 	rsp := user.NewUserResponseByLLM()
// 	t.Logf("ProcessChatCompletion response: %s", responseJSON)
// 	userResponse, err := rsp.UnmarshalToUserResponceByLLM(responseJSON)
// 	if err != nil {
// 		t.Fatalf("UnmarshalToUserResponceByLLM failed: %v", err)
// 	}

// 	// 验证解析结果
// 	if userResponse.PossibleNation == "" {
// 		t.Errorf("Expected non-empty PossibleNation, got empty string")
// 	}
// 	if userResponse.ConfidenceLevel == "" {
// 		t.Errorf("Expected non-empty ConfidenceLevel, got empty string")
// 	}

// }
