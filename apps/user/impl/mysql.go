package impl

import (
	"database/sql"
	"log"

	"github.com/acd19ml/TalentRank/apps"
	"github.com/acd19ml/TalentRank/apps/git"
	"github.com/acd19ml/TalentRank/apps/user"
	"github.com/acd19ml/TalentRank/conf"
	"google.golang.org/grpc"
)

type ServiceImpl struct {
	db   *sql.DB
	svc  git.GitServiceClient
	llm  user.LLMService
	rsp  user.LLMResponseService
	conn *grpc.ClientConn
}

var svcimpl = &ServiceImpl{}

func (s *ServiceImpl) Config() {
	// 配置数据库连接
	s.db = conf.C().MySQL.GetDB()

	// 配置 gRPC 连接
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}
	s.conn = conn // 保存连接以便稍后关闭
	s.svc = git.NewGitServiceClient(conn)
}

func (s *ServiceImpl) Name() string {
	return user.AppName
}

// _ import app 自动执行注册逻辑
func init() {
	//  对象注册到ioc层
	apps.RegistryImpl(svcimpl)
}
