package git_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/acd19ml/TalentRank/utils/git"

	"github.com/acd19ml/TalentRank/utils"

	"github.com/stretchr/testify/assert"
)

const (
	username = utils.Username
)

var client utils.Service

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

func TestGetRepoStars(t *testing.T) {
	ctx := context.Background()
	repoName := "wgan-gp" // 你可以替换为你想要测试的 GitHub 仓库名

	// 调用 GetRepoStars 方法
	stars, err := client.GetRepoStars(ctx, username, repoName)

	// 检查是否返回了错误
	assert.NoError(t, err)

	// 检查 stars 数量是否非负
	assert.GreaterOrEqual(t, stars, 0, "Stars count should be non-negative")
}

func TestGetCommits(t *testing.T) {

	ctx := context.Background()

	charLimit := 200 // 设定一个字符限制

	commits, err := client.GetCommits(ctx, username, charLimit)
	assert.NoError(t, err)
	assert.NotNil(t, commits)
	assert.LessOrEqual(t, len(commits), charLimit*21) // 假设最多100个repo
}

func TestGetFollowers(t *testing.T) {
	ctx := context.Background()
	followers, err := client.GetFollowers(ctx, username)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, followers, 0, "Followers count should be non-negative")
}

func TestGetTotalStars(t *testing.T) {
	ctx := context.Background()
	totalStars, err := client.GetTotalStars(ctx, username)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, totalStars, 0, "Total stars count should be non-negative")
}

func TestGetTotalForks(t *testing.T) {
	ctx := context.Background()
	totalForks, err := client.GetTotalForks(ctx, username)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, totalForks, 0, "Total forks count should be non-negative")
}

func TestGetRepositories(t *testing.T) {
	ctx := context.Background()
	repos, err := client.GetRepositories(ctx, username)
	assert.NoError(t, err)
	assert.NotEmpty(t, repos, "Repositories list should not be empty")
}

func TestGetRepoStarsMap(t *testing.T) {
	ctx := context.Background()
	// username := "octocat" // 使用 GitHub 上一个公开的示例用户

	// 调用 GetRepoStarsMap 方法
	repoStarsMap, err := client.GetRepoStarsMap(ctx, username)
	fmt.Println(repoStarsMap)

	// 检查是否返回了错误
	assert.NoError(t, err)

	// 检查返回的 map 是否非空
	assert.NotNil(t, repoStarsMap)
	assert.NotEmpty(t, repoStarsMap, "Repo stars map should not be empty")

	// 检查每个仓库的 star 数量是否非负
	for repo, stars := range repoStarsMap {
		assert.GreaterOrEqual(t, stars, 0, "Stars count for repo %s should be non-negative", repo)
	}
}

func TestGetOrganizations(t *testing.T) {
	ctx := context.Background()
	orgs, err := client.GetOrganizations(ctx, username)
	assert.NoError(t, err)
	assert.NotNil(t, orgs, "Organizations list should not be nil")
}
