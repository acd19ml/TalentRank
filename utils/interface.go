package utils

import "context"

type Service interface {
	// Rank
	GetFollowers(ctx context.Context, username string) (int, error)
	GetTotalStars(ctx context.Context, username string) (int, error)
	GetTotalForks(ctx context.Context, username string) (int, error)
	GetDependentRepositories(ctx context.Context, username string) (int, error)

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
	GetRepositories(ctx context.Context, username string) ([]string, error)
}
