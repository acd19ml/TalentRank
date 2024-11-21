package impl

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/acd19ml/TalentRank/apps"
	"github.com/acd19ml/TalentRank/apps/git"
	"github.com/acd19ml/TalentRank/apps/llm"
	"github.com/acd19ml/TalentRank/apps/user"
	"github.com/acd19ml/TalentRank/apps/user/kafka"
	"github.com/acd19ml/TalentRank/conf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ServiceImpl struct {
	Db       *sql.DB
	svc      git.GitServiceClient
	llm      llm.LLMServiceClient
	gitConn  *grpc.ClientConn
	llmConn  *grpc.ClientConn
	Producer user.MessageProducer
	Consumer user.MessageConsumer
}

var svcimpl = &ServiceImpl{}

func (s *ServiceImpl) Config() {
	// 配置数据库连接
	s.Db = conf.C().MySQL.GetDB()

	// 初始化 gRPC 连接
	if err := conf.C().App.InitClientConn(); err != nil {
		panic(err)
	}

	// 获取 LLM 服务连接
	llmConn, err := conf.C().App.GetServiceConn(user.LlmClient)
	if err != nil {
		log.Fatalf("failed to get LLM client service connection: %v", err)
	}

	// 获取 Git 服务连接
	gitConn, err := conf.C().App.GetServiceConn(user.GitClient)
	if err != nil {
		log.Fatalf("failed to get Git client service connection: %v", err)
	}

	s.llmConn = llmConn
	s.gitConn = gitConn

	// 初始化 gRPC 客户端
	s.svc = git.NewGitServiceClient(s.gitConn)
	s.llm = llm.NewLLMServiceClient(s.llmConn)

	brokers := []string{"localhost:9092"}

	s.InitializeKafka(brokers)

	// 初始化速率限制器
	rateLimiter := NewRateLimiter(10, time.Second) // 每秒最多处理 10 条消息

	// 启动生产者协程
	go s.StartProducer(context.Background())

	// 启动消费者协程
	go s.StartConsumer(context.Background(), rateLimiter)

}

// SetLLMClient 提供一个用于测试的 Setter 方法
func (s *ServiceImpl) SetLLMClient(client llm.LLMServiceClient) {
	s.llm = client
}

// SetGitClient 提供一个用于测试的 Setter 方法
func (s *ServiceImpl) SetGitClient(client git.GitServiceClient) {
	s.svc = client
}

func (s *ServiceImpl) NewAuthenticatedContext(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs())
}

func (s *ServiceImpl) Name() string {
	return user.AppName
}

func (s *ServiceImpl) InitializeKafka(brokers []string) error {
	log.Println("Initializing Kafka Producer...")
	s.Producer = kafka.NewKafkaProducer(brokers)

	log.Println("Initializing Kafka Consumers...")
	repoConsumer := kafka.NewKafkaConsumer(brokers, "repo_api_tasks", "repo-group")
	if repoConsumer == nil {
		return fmt.Errorf("failed to initialize Kafka consumer for topic: repo_api_tasks")
	}
	userConsumer := kafka.NewKafkaConsumer(brokers, "user_api_tasks", "user-group")
	if userConsumer == nil {
		return fmt.Errorf("failed to initialize Kafka consumer for topic: user_api_tasks")
	}

	// 将消费者注册到 ServiceImpl
	s.Consumer = &kafka.MultiTopicConsumer{
		RepoConsumer: repoConsumer,
		UserConsumer: userConsumer,
	}

	log.Println("Kafka Producer and Consumers initialized successfully.")
	return nil
}

// _ import app 自动执行注册逻辑
func init() {
	//  对象注册到ioc层
	apps.RegistryImpl(svcimpl)
}
