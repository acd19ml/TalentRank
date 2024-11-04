package impl

import (
	"context"
	"os"
	"sync"

	"github.com/acd19ml/TalentRank/apps"
	"github.com/acd19ml/TalentRank/apps/git"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

var svr = &Service{}

type Service struct {
	git.UnimplementedGitServiceServer
	client        *github.Client
	oauth         *oauth2.TokenSource
	reposCache    []string   // 缓存仓库列表
	cacheUsername string     // 缓存中保存的用户名
	cacheMutex    sync.Mutex // 用于同步缓存访问
}

// Config 配置服务
func (s *Service) Config() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	s.client = client
	s.oauth = &ts

}

func (s *Service) Name() string {
	return git.AppName
}

func (s *Service) Registry(server *grpc.Server) {
	git.RegisterGitServiceServer(server, svr)
}

func init() {
	apps.RegistryGrpc(svr)
}
