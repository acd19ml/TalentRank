package protocol

import (
	"fmt"
	"log"
	"net"

	"github.com/acd19ml/TalentRank/apps"
	"github.com/acd19ml/TalentRank/conf"
	"github.com/acd19ml/TalentRank/middleware/server"
	"google.golang.org/grpc"
)

func NewGRPCGitService() *GRPCService {
	// 添加认证中间件
	reqAuth := server.NewAuthUnaryServerInterceptor()
	svr := grpc.NewServer(
		grpc.UnaryInterceptor(reqAuth),
	)
	return &GRPCService{
		svr: svr,
		c:   conf.C(),
	}
}

type GRPCService struct {
	svr *grpc.Server
	c   *conf.Config
}

func (s *GRPCService) InitGRPC() {

	// 初始化所有grpc服务
	apps.InitGrpc(s.svr)

	apps.LoadedGrpcApps()
	log.Printf("loaded grpc apps: %v", apps.LoadedGrpcApps())
}

func (s *GRPCService) StartGit() {

	// 启动grpc服务
	lis, err := net.Listen("tcp", s.c.App.GitAddr())
	if err != nil {
		log.Printf("listen git grpc tcp conn error: %v", err)
		return
	}

	fmt.Printf("git grpc server start at %s\n", s.c.App.GitAddr())
	if err := s.svr.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			log.Printf("git grpc server stop")
		}

		log.Printf("git grpc server start error: %v", err.Error())
		return
	}
}

func (s *GRPCService) StartLlm() {

	// 启动grpc服务
	lis, err := net.Listen("tcp", s.c.App.LlmAddr())
	if err != nil {
		log.Printf("listen llm grpc tcp conn error: %v", err)
		return
	}

	fmt.Printf("llm grpc server start at %s\n", s.c.App.LlmAddr())
	if err := s.svr.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			log.Printf("llm grpc server stop")
		}

		log.Printf("llm grpc server start error: %v", err.Error())
		return
	}
}

func (s *GRPCService) Stop() error {
	s.svr.GracefulStop()
	return nil
}
