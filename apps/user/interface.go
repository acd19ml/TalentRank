package user

import (
	"context"
	"os/user"
)

type Service interface {
	CreateUser(context.Context, string) (*user.User, error)
}
