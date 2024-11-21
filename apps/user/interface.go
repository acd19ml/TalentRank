package user

import (
	"context"
)

type Service interface {
	CreateUserRepos(ctx context.Context, username string) (*UserRepos, error)
	QueryUsers(context.Context, *QueryUserRequest) (*UserSet, error)
	DescribeUserRepos(context.Context, *DescribeUserReposRequest) (string, error)
	GetLocationCounts(context.Context) ([]*GetLocationCountsRequest, error)
	DeleteUserRepos(context.Context, *DeleteUserReposRequest) (*DeleteUserReposResponse, error)
}

// 消息生产者接口
type MessageProducer interface {
	Produce(ctx context.Context, topic string, message interface{}) error
}

// 消息消费者接口
type MessageConsumer interface {
	Consume(ctx context.Context, topic string) ([]byte, error)
}
