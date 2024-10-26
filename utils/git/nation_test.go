package git_test

import (
	"acd19ml/TalentRank/utils"
	"acd19ml/TalentRank/utils/git"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	username = utils.Username
)

var client *git.Git

func init() {
	if os.Getenv("GITHUB_TOKEN") == "" {
		panic("GITHUB_TOKEN is not set. Please set it before running tests.")
	}
	client = git.NewGitClient() // 初始化客户端
}

func TestGetName(t *testing.T) {
	ctx := context.Background()
	name, err := client.GetName(ctx, username)
	assert.NoError(t, err)
	assert.NotEmpty(t, name)
}

func TestGetCompany(t *testing.T) {

	ctx := context.Background()

	company, err := client.GetCompany(ctx, username)
	assert.NoError(t, err)
	// 公司可能为空
	assert.NotNil(t, company)
}

func TestGetLocation(t *testing.T) {

	ctx := context.Background()

	location, err := client.GetLocation(ctx, username)
	assert.NoError(t, err)
	assert.NotNil(t, location)
}

func TestGetEmail(t *testing.T) {

	ctx := context.Background()

	email, err := client.GetEmail(ctx, username)
	assert.NoError(t, err)
	// 邮箱可能为空
	assert.NotNil(t, email)
}

func TestGetBio(t *testing.T) {

	ctx := context.Background()

	bio, err := client.GetBio(ctx, username)
	assert.NoError(t, err)
	// Bio 可能为空
	assert.NotNil(t, bio)
}

func TestGetReadme(t *testing.T) {
	ctx := context.Background()
	charLimit := 200 // 设定一个字符限制

	content, err := client.GetReadme(ctx, username, charLimit)
	assert.NoError(t, err)
	assert.NotNil(t, content)
	assert.LessOrEqual(t, len(content), charLimit*21) // 假设最多100个repo
}

func TestGetCommits(t *testing.T) {

	ctx := context.Background()

	charLimit := 200 // 设定一个字符限制

	commits, err := client.GetCommits(ctx, username, charLimit)
	assert.NoError(t, err)
	assert.NotNil(t, commits)
	assert.LessOrEqual(t, len(commits), charLimit*21) // 假设最多100个repo
}
