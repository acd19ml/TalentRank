package apps

import (
	"fmt"

	"github.com/acd19ml/TalentRank/apps/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

const (
	Username  = "lululxvi"
	CharLimit = 200
	RepoLimit = 10
)

var (
	UserService user.Service
	LlmService  user.LLMService
	implApps    = map[string]ImplService{}
	ginApps     = map[string]GinService{}
	grpcApps    = map[string]GrpcService{}
)

type GrpcService interface {
	Registry(server *grpc.Server)
	Config()
	Name() string
}

func RegistryGrpc(svc GrpcService) {
	// 检查服务是否已经注册过
	if _, ok := grpcApps[svc.Name()]; ok {
		// 如果该服务已经注册，抛出一个 panic 错误，避免重复注册
		panic(fmt.Sprintf("service %s already registered", svc.Name()))
	}

	// 将服务注册到 svcs 容器，键是服务的名字，值是该服务的实例
	grpcApps[svc.Name()] = svc
}

type ImplService interface {
	Config()
	Name() string
}

func RegistryImpl(svc ImplService) {
	// 检查服务是否已经注册过
	if _, ok := implApps[svc.Name()]; ok {
		// 如果该服务已经注册，抛出一个 panic 错误，避免重复注册
		panic(fmt.Sprintf("service %s already registered", svc.Name()))
	}

	// 将服务注册到 svcs 容器，键是服务的名字，值是该服务的实例
	implApps[svc.Name()] = svc
	// 判断传入的服务是否满足 form.Service 接口
	if v, ok := svc.(user.Service); ok {
		// 如果是 user.Service 类型的服务，将其赋值给全局的 UserService
		UserService = v
	}
}

// Get 一个Impl服务的实例：implApps
// 返回一个对象, 任何类型都可以, 使用时, 由使用方进行断言
func GetImpl(name string) interface{} {
	for k, v := range implApps {
		if k == name {
			return v
		}
	}

	return nil
}

type GinService interface {
	Registry(r gin.IRouter)
	Config()
	Name() string
}

func RegistryGin(svc GinService) {
	// 检查服务是否已经注册过
	if _, ok := ginApps[svc.Name()]; ok {
		// 如果该服务已经注册，抛出一个 panic 错误，避免重复注册
		panic(fmt.Sprintf("service %s already registered", svc.Name()))
	}

	// 将服务注册到 svcs 容器，键是服务的名字，值是该服务的实例
	ginApps[svc.Name()] = svc
}

// 已经加载完成的Gin App由哪些, 用于日志输出
func LoadedGinApps() (names []string) {
	for k := range ginApps {
		names = append(names, k)
	}
	return
}

func LoadedGrpcApps() (names []string) {
	for k := range grpcApps {
		names = append(names, k)
	}
	return
}

func InitGrpc(server *grpc.Server) {
	for _, v := range grpcApps {
		v.Config()
	}

	for _, v := range grpcApps {
		v.Registry(server)
	}
}

func InitGin(r gin.IRouter) {

	// 初始化对象
	for _, v := range ginApps {
		v.Config()
	}

	// 完成http handler的注册
	for _, v := range ginApps {
		v.Registry(r)
	}

}

// 用户初始化 注册到Ioc容器里面的所有服务
func InitImpl() {
	for _, v := range grpcApps {
		v.Config()
	}

	for _, v := range implApps {
		v.Config()
	}
}
