package user

import (
	"context"
)

type Service interface {
	CreateUserRepos(ctx context.Context, username string) (*UserRepos, error)
	QueryUserRepos(context.Context, *QueryUserReposRequest) (*UserReposSet, error)
}

type LLMService interface {
	ProcessChatCompletion(inputJSON []byte) ([]byte, error)
}
