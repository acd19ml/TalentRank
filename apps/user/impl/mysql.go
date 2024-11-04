package impl

import (
	"database/sql"
	"log"

	"github.com/acd19ml/TalentRank/apps"
	"github.com/acd19ml/TalentRank/apps/git"
	"github.com/acd19ml/TalentRank/apps/llm"
	"github.com/acd19ml/TalentRank/apps/user"
	"github.com/acd19ml/TalentRank/conf"
	"google.golang.org/grpc"
)

type ServiceImpl struct {
	db      *sql.DB
	svc     git.GitServiceClient
	llm     llm.LLMServiceClient
	gitConn *grpc.ClientConn
	llmConn *grpc.ClientConn
}

var svcimpl = &ServiceImpl{}

func (s *ServiceImpl) Config() {
	// 配置数据库连接
	s.db = conf.C().MySQL.GetDB()

	// 初始化 gRPC 连接
	s.gitConn = s.createGRPCConn("localhost:50051", "Git gRPC server")
	s.llmConn = s.createGRPCConn("localhost:50052", "LLM gRPC server")

	// 初始化 gRPC 客户端
	s.svc = git.NewGitServiceClient(s.gitConn)
	s.llm = llm.NewLLMServiceClient(s.llmConn)
}

// SetLLMClient 提供一个用于测试的 Setter 方法
func (s *ServiceImpl) SetLLMClient(client llm.LLMServiceClient) {
	s.llm = client
}

// SetGitClient 提供一个用于测试的 Setter 方法
func (s *ServiceImpl) SetGitClient(client git.GitServiceClient) {
	s.svc = client
}

// createGRPCConn 创建 gRPC 连接并返回连接对象
func (s *ServiceImpl) createGRPCConn(address string, serviceName string) *grpc.ClientConn {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to %s: %v", serviceName, err)
	}
	return conn
}

func (s *ServiceImpl) Name() string {
	return user.AppName
}

// _ import app 自动执行注册逻辑
func init() {
	//  对象注册到ioc层
	apps.RegistryImpl(svcimpl)
}
