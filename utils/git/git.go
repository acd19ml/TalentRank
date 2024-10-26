package git

import (
	"acd19ml/TalentRank/utils"
	"context"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// 接口检查
var _ utils.Rank = (*git)(nil)

func NewGitClient() *git {
	ctx := context.Background()
	// 使用 OAuth2 令牌进行认证
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &git{
		client: client,
		ctx:    ctx,
		oauth:  &ts,
	}
}

type git struct {
	client *github.Client
	ctx    context.Context
	oauth  *oauth2.TokenSource
}

func (g *git) GetUser(username string) (*github.User, error) {
	user, _, err := g.client.Users.Get(g.ctx, username)
	if err != nil {
		panic(err)
	}
	return user, nil
}
