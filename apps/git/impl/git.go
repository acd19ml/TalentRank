package impl

import (
	"context"
	"log"
	"os"
	_ "os"
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
	// client        *github.Client
	reposCache    []string   // 缓存仓库列表
	cacheUsername string     // 缓存中保存的用户名
	cacheMutex    sync.Mutex // 用于同步缓存访问
}

// Config 配置服务
func (s *Service) Config() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is not set")
	}
	// s.client = s.GetClientWithToken(token)

}

//// Config 配置服务
//func (s *Service) Config(token string) {
//	if token == "" {
//		log.Fatal("GITHUB_TOKEN is not provided")
//	}
//
//	// 使用传入的 Token 配置 OAuth2 客户端
//	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
//	tc := oauth2.NewClient(context.Background(), ts)
//	client := github.NewClient(tc)
//
//	// 将配置应用到服务实例
//	s.client = client
//	s.oauth = &ts
//}

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
