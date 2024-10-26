package utils

import "context"

type Rank interface {
	GetFollowers(ctx context.Context, username string) (int, error)
	GetTotalStars(ctx context.Context, username string) (int, error)
	GetTotalForks(ctx context.Context, username string) (int, error)
	GetOrganizations(ctx context.Context, username string) ([]string, error)
}

type Nation interface {
}
