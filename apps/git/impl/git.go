package impl

import (
	"context"
	"log"
	"os"
	_ "os"

	"github.com/acd19ml/TalentRank/apps"
	"github.com/acd19ml/TalentRank/apps/git"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

var svr = &Service{}

type Service struct {
	git.UnimplementedGitServiceServer
}

// Config 配置服务
func (s *Service) Config() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is not set")
	}
}

func (s *Service) Name() string {
	return git.AppName
}

func (s *Service) Registry(server *grpc.Server) {
	git.RegisterGitServiceServer(server, svr)
}

func (s *Service) GetClientWithToken(token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}

// 获取 client 时使用上下文token
func (s *Service) getClientFromContext(ctx context.Context) *github.Client {
	token, ok := ctx.Value("githubToken").(string)
	if !ok || token == "" {
		defaultToken := os.Getenv("GITHUB_TOKEN")
		return s.GetClientWithToken(defaultToken)
	}
	return s.GetClientWithToken(token)
}

func init() {
	apps.RegistryGrpc(svr)
}
