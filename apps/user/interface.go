package user

import (
	"context"
)

type Service interface {
	CreateUserRepos(ctx context.Context, username string) (*UserRepos, error)
	QueryUsers(context.Context, *QueryUserRequest) (*UserSet, error)
	DescribeUserRepos(context.Context, *DescribeUserReposRequest) (string, error)
}
