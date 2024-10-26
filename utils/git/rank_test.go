package git_test

import (
	"acd19ml/TalentRank/utils/git"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	if os.Getenv("GITHUB_TOKEN") == "" {
		panic("GITHUB_TOKEN is not set. Please set it before running tests.")
	}
	client = git.NewGitClient() // 初始化客户端
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

func TestGetOrganizations(t *testing.T) {
	ctx := context.Background()
	orgs, err := client.GetOrganizations(ctx, username)
	assert.NoError(t, err)
	assert.NotNil(t, orgs, "Organizations list should not be nil")
}
