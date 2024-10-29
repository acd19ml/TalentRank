package user

import (
	"context"
)

type Service interface {
	CreateUserRepos(ctx context.Context, username string) (*UserRepos, error)
	QueryUserRepos(context.Context, *QueryUserReposRequest) (*UserReposSet, error)
	// GetUserReposJSON(ctx context.Context, username string) (string, error)
}

type LLMService interface {
	// GetUserReposJSON(ctx context.Context, username string) (string, error)
	GetUserReposJSONWithRequest(ctx context.Context, username string) (string, error)
}
