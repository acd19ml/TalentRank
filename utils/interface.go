package utils

import "context"

type Service interface {
	// Rank
	GetFollowers(ctx context.Context, username string) (int, error)
	GetTotalStars(ctx context.Context, username string) (int, error)
	GetTotalForks(ctx context.Context, username string) (int, error)
	GetDependentRepositories(ctx context.Context, username string) (int, error)
	GetStarsByRepo(ctx context.Context, username string) (map[string]int, error)
	GetTotalCommitsByRepo(ctx context.Context, username string) (map[string]int, error)
	GetUserCommitsByRepo(ctx context.Context, username string) (map[string]int, error)
	GetTotalIssuesByRepo(ctx context.Context, username string) (map[string]int, error)
	GetUserSolvedIssuesByRepo(ctx context.Context, username string) (map[string]int, error)
	GetTotalPullRequestsByRepo(ctx context.Context, username string) (map[string]int, error)
	GetUserMergedPullRequestsByRepo(ctx context.Context, username string) (map[string]int, error)
	GetTotalCodeReviewsByRepo(ctx context.Context, username string) (map[string]int, error)
	GetUserCodeReviewsByRepo(ctx context.Context, username string) (map[string]int, error)
	// Nation
	GetName(ctx context.Context, username string) (string, error)
	GetCompany(ctx context.Context, username string) (string, error)
	GetLocation(ctx context.Context, username string) (string, error)
	GetEmail(ctx context.Context, username string) (string, error)
	GetBio(ctx context.Context, username string) (string, error)
	GetOrganizations(ctx context.Context, username string) ([]string, error)
	GetReadme(ctx context.Context, username string, charLimit int) (string, error)
	GetCommits(ctx context.Context, username string, charLimit int) (string, error)
	// Tools
	GetRepoStars(ctx context.Context, owner, repoName string) (int, error)
	GetRepositories(ctx context.Context, username string) ([]string, error)
	GetTotalLineChanges(ctx context.Context, username string) (int, error)
	GetLineChanges(ctx context.Context, username, repoName string) (int, error)
}
