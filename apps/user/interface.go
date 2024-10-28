package user

import (
	"context"
)

type Service interface {
	CreateUser(context.Context, string) (*User, error)
	QueryUser(context.Context, *QueryUserReposRequest) (*UserReposSet, error)
}
