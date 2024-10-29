package impl

import (
	"context"
	"fmt"

	"github.com/acd19ml/TalentRank/apps/user"
)

func (s *ServiceImpl) CreateUserRepos(ctx context.Context, username string) (*user.UserRepos, error) {
	ins, err := s.constructUserRepos(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to construct user repos: %w", err)
	}

	if err = s.save(ctx, ins); err != nil {
		return nil, fmt.Errorf("failed to save user repos: %w", err)
	}
	return ins, nil
}

func (s *ServiceImpl) QueryUserRepos(ctx context.Context, req *user.QueryUserReposRequest) (*user.UserReposSet, error) {

	return nil, nil
}
