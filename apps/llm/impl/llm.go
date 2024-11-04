package impl

import (
	"os"

	"github.com/acd19ml/TalentRank/apps"
	"github.com/acd19ml/TalentRank/apps/llm"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"google.golang.org/grpc"
)

var svr = &LLMServer{}

type LLMServer struct {
	llm.UnimplementedLLMServiceServer
	client *arkruntime.Client
}

func (s *LLMServer) Config() {
	s.client = arkruntime.NewClientWithApiKey(
		os.Getenv("ARK_API_KEY"),
		arkruntime.WithBaseUrl("https://ark.cn-beijing.volces.com/api/v3"),
		arkruntime.WithRegion("cn-beijing"),
	)
}

func (s *LLMServer) Name() string {
	return llm.AppName
}

func (s *LLMServer) Registry(server *grpc.Server) {
	llm.RegisterLLMServiceServer(server, s)
}

func init() {
	apps.RegistryGrpc(svr)
}
