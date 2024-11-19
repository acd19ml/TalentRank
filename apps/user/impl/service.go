package impl

import (
	"context"
	"database/sql"

	"github.com/acd19ml/TalentRank/apps"
	"github.com/acd19ml/TalentRank/apps/git"
	"github.com/acd19ml/TalentRank/apps/llm"
	"github.com/acd19ml/TalentRank/apps/user"
	"github.com/acd19ml/TalentRank/conf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ServiceImpl struct {
	Db      *sql.DB
	svc     git.GitServiceClient
	llm     llm.LLMServiceClient
	gitConn *grpc.ClientConn
	llmConn *grpc.ClientConn
}

var svcimpl = &ServiceImpl{}

func (s *ServiceImpl) Config() {
	// 配置数据库连接
	s.Db = conf.C().MySQL.GetDB()

	// 初始化 gRPC 连接
	s.gitConn = conf.C().GRPCPool.InitClientConn(user.GitClient)

	s.llmConn = conf.C().GRPCPool.InitClientConn(user.LlmClient)

	// 初始化 gRPC 客户端
	s.svc = git.NewGitServiceClient(s.gitConn)
	s.llm = llm.NewLLMServiceClient(s.llmConn)

	// 启动定时任务
	go func() {
		s.StartWeeklyUpdate(s.NewAuthenticatedContext(), apps.UpdateInterval)
	}()
}

// SetLLMClient 提供一个用于测试的 Setter 方法
func (s *ServiceImpl) SetLLMClient(client llm.LLMServiceClient) {
	s.llm = client
}

// SetGitClient 提供一个用于测试的 Setter 方法
func (s *ServiceImpl) SetGitClient(client git.GitServiceClient) {
	s.svc = client
}

func (s *ServiceImpl) NewAuthenticatedContext() context.Context {
	return metadata.NewOutgoingContext(context.Background(), metadata.Pairs())
}

func (s *ServiceImpl) Name() string {
	return user.AppName
}

// _ import app 自动执行注册逻辑
func init() {
	//  对象注册到ioc层
	apps.RegistryImpl(svcimpl)
}
