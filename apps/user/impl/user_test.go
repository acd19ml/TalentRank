package impl_test

import (
	"context"
	"os"
	"testing"

	"github.com/acd19ml/TalentRank/apps/git"
	"github.com/acd19ml/TalentRank/apps/user"
	"github.com/acd19ml/TalentRank/apps/user/impl"
	"github.com/acd19ml/TalentRank/conf"
)

const (
	username = git.Username
)

var (
	client user.Service
	ctx    = context.Background()
)

func init() {
	if os.Getenv("GITHUB_TOKEN") == "" {
		panic("GITHUB_TOKEN is not set. Please set it before running tests.")
	}

	err := conf.LoadConfigFromToml("../../../etc/demo.toml")
	if err != nil {
		panic(err)
	}

	client = impl.NewUserServiceImpl() // 初始化客户端

}

func TestCreateUserRepos(t *testing.T) {

	// 调用 CreateUserRepos
	ctx := context.Background()
	userRepos, err := client.CreateUserRepos(ctx, username)

	// 检查是否返回了错误
	if err != nil {
		t.Fatalf("CreateUserRepos returned an error: %v", err)
	}

	// 校验返回的结果
	if userRepos == nil {
		t.Fatal("Expected non-nil UserRepos, got nil")
	}

	// 这里可以进行更具体的检查，比如检查 userRepos 的字段值是否符合预期
	if userRepos.Username != username {
		t.Errorf("Expected Username to be %s, got %s", username, userRepos.Username)
	}
}
