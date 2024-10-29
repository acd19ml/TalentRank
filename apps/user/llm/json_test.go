package llm_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/acd19ml/TalentRank/apps/git"
	"github.com/acd19ml/TalentRank/apps/user"
	"github.com/acd19ml/TalentRank/apps/user/llm"
	"github.com/acd19ml/TalentRank/conf"
	"github.com/stretchr/testify/assert"
)

const (
	username = git.Username
)

var (
	ctx    = context.Background()
	client user.LLMService
)

func init() {
	if os.Getenv("GITHUB_TOKEN") == "" {
		panic("GITHUB_TOKEN is not set. Please set it before running tests.")
	}

	err := conf.LoadConfigFromToml("../../../etc/demo.toml")
	if err != nil {
		panic(err)
	}

	client = llm.NewLLMJsonHandler() // 初始化客户端

}

func TestLoadConfigFromToml(t *testing.T) {
	should := assert.New(t)
	err := conf.LoadConfigFromToml("../../../etc/demo.toml")
	if should.NoError(err) {
		should.Equal("TalentRank", conf.C().App.Name)
	}
}

func TestGetUserReposJSONWithRequest(t *testing.T) {

	// 调用 GetUserReposJSONWithRequest
	result, err := client.GetUserReposJSONWithRequest(ctx, username)

	// 检查是否出错
	assert.NoError(t, err)

	// 检查结果是否为有效的 JSON
	var parsedResult map[string]interface{}
	err = json.Unmarshal([]byte(result), &parsedResult)
	assert.NoError(t, err, "Result should be a valid JSON")

	t.Logf("Result: %s", result)
}
