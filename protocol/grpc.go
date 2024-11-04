package protocol

import (
	"fmt"
	"log"
	"net"

	"github.com/acd19ml/TalentRank/apps"
	"github.com/acd19ml/TalentRank/conf"
	"google.golang.org/grpc"
)

func NewGRPCService() *GRPCService {
	svr := grpc.NewServer()
	return &GRPCService{
		svr: svr,
		c:   conf.C(),
	}
}

type GRPCService struct {
	svr *grpc.Server
	c   *conf.Config
}

func (s *GRPCService) Start() {
	// 初始化所有grpc服务
	apps.InitGrpc(s.svr)

	apps.LoadedGrpcApps()
	log.Printf("loaded grpc apps: %v", apps.LoadedGrpcApps())

	// 启动grpc服务
	lis, err := net.Listen("tcp", s.c.App.GrpcAddr())
	if err != nil {
		log.Printf("listen grpc tcp conn error: %v", err)
		return
	}

	fmt.Printf("grpc server start at %s", s.c.App.GrpcAddr())
	if err := s.svr.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			log.Printf("grpc server stop")
		}

		log.Printf("grpc server start error: %v", err.Error())
		return
	}
}

func (s *GRPCService) Stop() error {
	s.svr.GracefulStop()
	return nil
}
