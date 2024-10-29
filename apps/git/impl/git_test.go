package impl_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/acd19ml/TalentRank/apps"
	"github.com/acd19ml/TalentRank/apps/git"
	"github.com/acd19ml/TalentRank/apps/git/impl"
	"github.com/stretchr/testify/assert"
)

const (
	username = apps.Username
)

var (
	client git.Service
	ctx    = context.Background()
)

func init() {
	if os.Getenv("GITHUB_TOKEN") == "" {
		panic("GITHUB_TOKEN is not set. Please set it before running tests.")
	}
	client = impl.NewGitClient() // 初始化客户端
}

func TestGetName(t *testing.T) {

	name, err := client.GetName(ctx, username)
	assert.NoError(t, err)
	assert.NotEmpty(t, name)
}

func TestGetCompany(t *testing.T) {

	company, err := client.GetCompany(ctx, username)
	assert.NoError(t, err)
	// 公司可能为空
	assert.NotNil(t, company)
}

func TestGetLocation(t *testing.T) {

	location, err := client.GetLocation(ctx, username)
	assert.NoError(t, err)
	assert.NotNil(t, location)
}

func TestGetEmail(t *testing.T) {

	email, err := client.GetEmail(ctx, username)
	assert.NoError(t, err)
	// 邮箱可能为空
	assert.NotNil(t, email)
}

func TestGetBio(t *testing.T) {

	bio, err := client.GetBio(ctx, username)
	assert.NoError(t, err)
	// Bio 可能为空
	assert.NotNil(t, bio)
}

func TestGetReadme(t *testing.T) {

	content, err := client.GetReadme(ctx, username, apps.CharLimit, apps.RepoLimit)
	assert.NoError(t, err)
	assert.NotNil(t, content)
	assert.LessOrEqual(t, len(content), apps.CharLimit+1) // 假设最多100个repo
}

func TestGetCommits(t *testing.T) {

	commits, err := client.GetCommits(ctx, username, apps.CharLimit, apps.RepoLimit)
	assert.NoError(t, err)
	assert.NotNil(t, commits)
	assert.LessOrEqual(t, len(commits), apps.CharLimit+1) // 假设最多100个repo
}

func TestGetFollowers(t *testing.T) {

	followers, err := client.GetFollowers(ctx, username)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, followers, 0, "Followers count should be non-negative")
}

func TestGetStarsByRepo(t *testing.T) {

	// username := "octocat" // 使用 GitHub 上一个公开的示例用户

	// 调用 GetStarsByRepo 方法
	repoStarsMap, err := client.GetStarsByRepo(ctx, username)

	// 检查是否返回了错误
	assert.NoError(t, err)

	// 检查返回的 map 是否非空
	assert.NotNil(t, repoStarsMap)
	assert.NotEmpty(t, repoStarsMap, "Repo stars map should not be empty")

	// 检查每个仓库的 star 数量是否非负
	for repo, stars := range repoStarsMap {
		assert.GreaterOrEqual(t, stars, 0, "Stars count for repo %s should be non-negative", repo)
		// t.Logf("Stars in repo %s: %d", repo, stars)
	}
}

func TestGetForksByRepo(t *testing.T) {
	ctx := context.Background()

	// 调用 GetForksByRepo 方法
	repoForksMap, err := client.GetForksByRepo(ctx, username)
	// 检查是否返回了错误
	assert.NoError(t, err)

	// 检查返回的 map 是否非空
	assert.NotNil(t, repoForksMap)
	assert.NotEmpty(t, repoForksMap, "Repo forks map should not be empty")

	// 检查每个仓库的 forks 数量是否非负
	for repo, forks := range repoForksMap {
		assert.GreaterOrEqual(t, forks, 0, "Forks count for repo %s should be non-negative", repo)
	}
}

func TestGetOrganizations(t *testing.T) {

	_, err := client.GetOrganizations(ctx, username)
	assert.NoError(t, err)
}

// func TestGetTotalCommitsByRepo(t *testing.T) {

// 	// 获取所有仓库的提交总数
// 	totalCommitsByRepo, err := client.GetTotalCommitsByRepo(ctx, username)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, totalCommitsByRepo, "Total commits map should not be nil")

// 	// 检查每个仓库的提交数是否非负
// 	for repo, commitCount := range totalCommitsByRepo {
// 		// t.Logf("Repo: %s, Total Commits: %d", repo, commitCount)
// 		assert.GreaterOrEqual(t, commitCount, 0, "Commit count should be non-negative for repository "+repo)
// 	}
// }

// func TestGetUserCommitsByRepo(t *testing.T) {

// 	// 获取用户在每个仓库的提交数
// 	userCommitsByRepo, err := client.GetUserCommitsByRepo(ctx, username)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, userCommitsByRepo, "User commits map should not be nil")

// 	// 检查每个仓库的用户提交数是否非负
// 	for repo, commitCount := range userCommitsByRepo {
// 		// t.Logf("Repo: %s, User Commits: %d", repo, commitCount)
// 		assert.GreaterOrEqual(t, commitCount, 0, "User commit count should be non-negative for repository "+repo)
// 	}
// }

func TestGetTotalIssuesByRepo(t *testing.T) {

	result, err := client.GetTotalIssuesByRepo(ctx, username)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	// for repo, count := range result {
	// 	t.Logf("Total issues for repo %s: %d", repo, count)
	// }
}

func TestGetUserSolvedIssuesByRepo(t *testing.T) {

	result, err := client.GetUserSolvedIssuesByRepo(ctx, username)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	// for repo, count := range result {
	// 	t.Logf("Solved issues by user in repo %s: %d", repo, count)
	// }
}

func TestGetTotalPullRequestsByRepo(t *testing.T) {

	result, err := client.GetTotalPullRequestsByRepo(ctx, username)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	// for repo, count := range result {
	// 	t.Logf("Total pull requests for repo %s: %d", repo, count)
	// }
}

func TestGetUserMergedPullRequestsByRepo(t *testing.T) {

	result, err := client.GetUserMergedPullRequestsByRepo(ctx, username)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	// for repo, count := range result {
	// 	t.Logf("Merged pull requests by user in repo %s: %d", repo, count)
	// }
}

func TestGetTotalCodeReviewsByRepo(t *testing.T) {

	result, err := client.GetTotalCodeReviewsByRepo(ctx, username)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	// for repo, count := range result {
	// 	t.Logf("Total code reviews for repo %s: %d", repo, count)
	// }
}

func TestGetUserCodeReviewsByRepo(t *testing.T) {

	result, err := client.GetUserCodeReviewsByRepo(ctx, username)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	// for repo, count := range result {
	// 	t.Logf("User's code reviews for repo %s: %d", repo, count)
	// }
}

func TestGetDependentRepositoriesByRepo(t *testing.T) {

	repoDependentsCount, err := client.GetDependentRepositoriesByRepo(ctx, username)

	// 验证结果
	assert.NoError(t, err)
	assert.NotEmpty(t, repoDependentsCount, "repoDependentsCount should not be empty")

	// 打印每个仓库的依赖数量
	for _, dependents := range repoDependentsCount {
		// t.Logf("Repo: %s, Dependents: %d", repo, dependents)
		assert.GreaterOrEqual(t, dependents, 0, "Dependents count should be non-negative")
	}
}

func TestGetLineChangesByRepo(t *testing.T) {

	// 调用 GetLineChangesByRepo 函数
	lineChanges, err := client.GetLineChangesByRepo(ctx, username)
	if err != nil {
		t.Fatalf("Failed to get line changes: %v", err)
	}

	// 输出结果并验证数据
	for repo, changes := range lineChanges {
		fmt.Printf("仓库: %s, 所有贡献者增删行数: %d, 用户 %s 的增删行数: %d\n", repo, changes[0], username, changes[1])

		// 简单验证增删行数
		if changes[0] < 0 || changes[1] < 0 {
			t.Errorf("Expected non-negative line changes, got %d for all contributors and %d for user %s in repo %s", changes[0], changes[1], username, repo)
		}
	}
}
