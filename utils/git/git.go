package git

import (
	"context"
	"os"

	"github.com/acd19ml/TalentRank/utils"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// 接口检查
var _ utils.Service = (*Git)(nil)

func NewGitClient() *Git {
	ctx := context.Background()
	// 使用 OAuth2 令牌进行认证
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &Git{
		client: client,
		oauth:  &ts,
	}
}

type Git struct {
	client *github.Client
	oauth  *oauth2.TokenSource
}

func (g *Git) GetUser(ctx context.Context, username string) (*github.User, error) {
	user, _, err := g.client.Users.Get(ctx, username)
	if err != nil {
		panic(err)
	}
	return user, nil
}

func (g *Git) GetRepositorie(ctx context.Context, username string, repo string) (*github.Repository, error) {
	repository, _, err := g.client.Repositories.Get(ctx, username, repo)
	if err != nil {
		panic(err)
	}
	return repository, nil
}
