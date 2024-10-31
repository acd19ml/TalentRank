package user

import (
	"context"
)

type Service interface {
	CreateUserRepos(ctx context.Context, username string) (*UserRepos, error)
	QueryUsers(context.Context, *QueryUserRequest) (*UserSet, error)
	DescribeUserRepos(context.Context, *DescribeUserReposRequest) (*UserRepos, error)
}

type LLMService interface {
	ProcessChatCompletion(inputJSON []byte) ([]byte, error)
}

type LLMResponseService interface {
	UnmarshalToUserResponceByLLM(data []byte) (*UserResponceByLLM, error)
}
