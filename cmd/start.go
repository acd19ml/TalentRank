package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/acd19ml/TalentRank/apps"
	_ "github.com/acd19ml/TalentRank/apps/all"
	"github.com/acd19ml/TalentRank/conf"
	"github.com/acd19ml/TalentRank/protocol"
	"github.com/spf13/cobra"
)

var (
	// pusher service config options
	confFile string
)

// 程序启动组装
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "启动TalentRank 后端API",
	Long:  "启动TalentRank 后端API",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 加载配置
		err := conf.LoadConfigFromToml(confFile)
		if err != nil {
			return (err)
		}

		// 初始化所有的服务
		apps.InitImpl()

		svc := newManager()

		ch := make(chan os.Signal, 1)
		// 中断信号处理
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGINT)
		// 协程服务不断监听信号
		go svc.WaitStop(ch)

		// 启动grpc服务
		svc.grpc.InitGRPC()
		go svc.grpc.StartGit()
		go svc.grpc.StartLlm()

		// 启动HTTP服务
		return svc.Start()
	},
}

func newManager() *manager {
	return &manager{
		http: protocol.NewHttpService(),
		grpc: protocol.NewGRPCGitService(),
	}
}

// 用于管理所有需要启动的服务
// 1. HTTP服务的启动
type manager struct {
	http *protocol.HttpService
	grpc *protocol.GRPCService
}

func (m *manager) Start() error {
	// 启动HTTP服务
	if err := m.http.Start(); err != nil {
		return err
	}
	return nil
}

// 处理来自外部的中断信号，Terminal信号
func (m *manager) Stop() error {
	return nil
}

// 处理来自外部的中断信号, 比如Terminal
func (m *manager) WaitStop(ch <-chan os.Signal) { //只读channel
	for v := range ch {
		switch v {
		default:
			// 先关闭内部
			if err := m.grpc.Stop(); err != nil {
				log.Printf("grpc stop error: %s", err)
			}
			log.Printf("receive signal: %s", v)
			m.http.Stop()
		}
	}

}

func init() {
	StartCmd.PersistentFlags().StringVarP(&confFile, "config", "f", "etc/demo.toml", "TalentRank api 配置文件路径")
	RootCmd.AddCommand(StartCmd)
}
