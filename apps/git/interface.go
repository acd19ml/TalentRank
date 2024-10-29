package git

import "context"

type Service interface {
	// Repo table

	GetDependentRepositoriesByRepo(ctx context.Context, username string) (map[string]int, error)
	GetStarsByRepo(ctx context.Context, username string) (map[string]int, error)
	GetForksByRepo(ctx context.Context, username string) (map[string]int, error)
	GetTotalIssuesByRepo(ctx context.Context, username string) (map[string]int, error)
	GetUserSolvedIssuesByRepo(ctx context.Context, username string) (map[string]int, error)
	GetTotalPullRequestsByRepo(ctx context.Context, username string) (map[string]int, error)
	GetUserMergedPullRequestsByRepo(ctx context.Context, username string) (map[string]int, error)
	GetTotalCodeReviewsByRepo(ctx context.Context, username string) (map[string]int, error)
	GetUserCodeReviewsByRepo(ctx context.Context, username string) (map[string]int, error)
	GetLineChangesByRepo(ctx context.Context, username string) (map[string][]int, error) //[]int{totalChanges, userChanges, totalCommits, userCommits}

	// User table
	GetName(ctx context.Context, username string) (string, error)
	GetCompany(ctx context.Context, username string) (string, error)
	GetLocation(ctx context.Context, username string) (string, error)
	GetEmail(ctx context.Context, username string) (string, error)
	GetBio(ctx context.Context, username string) (string, error)
	GetOrganizations(ctx context.Context, username string) ([]string, error)
	GetFollowers(ctx context.Context, username string) (int, error)
	GetReadme(ctx context.Context, username string, charLimit int) (string, error)
	GetCommits(ctx context.Context, username string, charLimit int) (string, error)
}
