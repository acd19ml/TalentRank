package impl_test

// import (
// 	"context"
// 	"os"
// 	"testing"

// 	"github.com/acd19ml/TalentRank/apps"
// 	"github.com/acd19ml/TalentRank/apps/user"
// 	"github.com/acd19ml/TalentRank/apps/user/impl"
// 	"github.com/acd19ml/TalentRank/apps/user/llm"
// 	_ "github.com/acd19ml/TalentRank/apps/user/llm"
// 	"github.com/acd19ml/TalentRank/conf"
// )

// const (
// 	username = apps.Username
// )

// var (
// 	client user.Service
// 	ctx    = context.Background()
// )

// func init() {
// 	if os.Getenv("GITHUB_TOKEN") == "" {
// 		panic("GITHUB_TOKEN is not set. Please set it before running tests.")
// 	}

// 	err := conf.LoadConfigFromToml("../../../etc/demo.toml")
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = llm.Init()
// 	if err != nil {
// 		panic("Failed to initialize ChatService")
// 	}

// 	client = impl.NewUserServiceImpl() // 初始化客户端

// }

// func TestCreateUserRepos(t *testing.T) {

// 	// 调用 CreateUserRepos
// 	ctx := context.Background()
// 	userRepos, err := client.CreateUserRepos(ctx, username)

// 	// 检查是否返回了错误
// 	if err != nil {
// 		t.Fatalf("CreateUserRepos returned an error: %v", err)
// 	}

// 	// 校验返回的结果
// 	if userRepos == nil {
// 		t.Fatal("Expected non-nil UserRepos, got nil")
// 	}

// 	// 校验返回的结果
// 	if userRepos.Username != username {
// 		t.Errorf("Expected Username to be %s, got %s", username, userRepos.Username)
// 	}
// }
