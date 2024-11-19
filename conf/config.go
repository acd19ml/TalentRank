package conf

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/acd19ml/TalentRank/middleware/client"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

// 全局config实例对象,
// 也就是我们程序，在内存中的配置对象
// 程序内部获取配置, 都通过读取该对象
// 该Config对象 什么时候被初始化?
//
//	配置加载时:
//	   LoadConfigFromToml
//	   LoadConfigFromEnv
//
// 为了不被程序在运行时恶意修改, 设置成私有变量
var config *Config

// 全局MySQL客户端实例
var db *sql.DB

// 要想获取配置, 单独提供函数
// 全局Config对象获取函数
func C() *Config {
	return config
}

// 初始化一个有默认值的Config对象
func NewDefaultConfig() *Config {
	return &Config{
		App:      NewDefaultApp(),
		MySQL:    NewDefaultMySQL(),
		GRPCPool: NewDefaultGRPCPool(),
	}
}

// Config 应用配置
// 通过封装为一个对象, 来与外部配置进行对接
type Config struct {
	App      *App      `toml:"app"`
	MySQL    *MySQL    `toml:"mysql"`
	GRPCPool *GRPCPool `toml:"grpc_pool"`
}

func NewDefaultApp() *App {
	return &App{
		Name: "demo",
		Host: "127.0.0.1",
		Port: "8050",
	}
}

type App struct {
	Name    string `toml:"name" env:"APP_NAME"`
	Host    string `toml:"host" env:"APP_HOST"`
	Port    string `toml:"port" env:"APP_PORT"`
	Key     string `toml:"key" env:"APP_KEY"`
	GitPort string `toml:"git_port" env:"APP_GIT_PORT"`
	LlmPort string `toml:"llm_port" env:"APP_LLM_PORT"`
}

func NewDefaultMySQL() *MySQL {
	return &MySQL{
		Host:        "127.0.0.1",
		Port:        "3306",
		UserName:    "root",
		Password:    "123456",
		Database:    "demo",
		MaxOpenConn: 200,
		MaxIdleConn: 50,
		MaxLifeTime: 1800,
		MaxIdleTime: 600,
	}
}

type MySQL struct {
	Host     string `toml:"host" env:"MYSQL_HOST"`
	Port     string `toml:"port" env:"MYSQL_PORT"`
	UserName string `toml:"username" env:"MYSQL_USERNAME"`
	Password string `toml:"password" env:"MYSQL_PASSWORD"`
	Database string `toml:"database" env:"MYSQL_DATABASE"`
	// 因为使用的MySQL连接池, 需要池做一些规划配置
	// 控制当前程序的MySQL打开的连接数
	MaxOpenConn int `toml:"max_open_conn" env:"MYSQL_MAX_OPEN_CONN"`
	// 控制MySQL复用, 比如5, 最多运行5个来复用
	MaxIdleConn int `toml:"max_idle_conn" env:"MYSQL_MAX_IDLE_CONN"`
	// 一个连接的生命周期, 这个和MySQL Server配置有关系, 必须小于Server配置
	// 一个连接用12h 换一个conn, 保证一定的可用性
	MaxLifeTime int `toml:"max_life_time" env:"MYSQL_MAX_LIFE_TIME"`
	// Idle 连接 最多允许存活多久
	MaxIdleTime int `toml:"max_idle_time" env:"MYSQL_MAX_idle_TIME"`
	// 作为私有变量, 用户与控制GetDB
	lock sync.Mutex
}

// 1. 第一种方式, 使用LoadGlobal 在加载时 初始化全局db实例
// 2. 第二种方式, 惰性加载, 获取DB是，动态判断再初始化
func (m *MySQL) GetDB() *sql.DB {
	// 直接加锁，锁住临界区
	m.lock.Lock()
	defer m.lock.Unlock()

	// 如果实例不存在，初始化一个新实例
	if db == nil {
		conn, err := m.getDBConn()
		if err != nil {
			panic(err)
		}
		db = conn
	}
	return db
}

// 连接池, driverConn具体的连接对象, 他维护着一个Socket
// pool []*driverConn, 维护pool里面的连接都是可用的, 定期检查我们的conn健康情况
// 某一个driverConn已经失效, driverConn.Reset(), 清空该结构体的数据, Reconn获取一个连接, 让该conn借壳存活
// 避免driverConn结构体的内存申请和释放的一个成本
func (m *MySQL) getDBConn() (*sql.DB, error) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&multiStatements=true", m.UserName, m.Password, m.Host, m.Port, m.Database)
	// 打开一个MySQL连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to mysql<%s> error, %s", dsn, err.Error())
	}
	db.SetMaxOpenConns(m.MaxOpenConn)
	db.SetMaxIdleConns(m.MaxIdleConn)
	db.SetConnMaxLifetime(time.Second * time.Duration(m.MaxLifeTime))
	db.SetConnMaxIdleTime(time.Second * time.Duration(m.MaxIdleTime))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping mysql<%s> error, %s", dsn, err.Error())
	}
	return db, nil
}

func (a *App) HTTPAddr() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}

func (a *App) GitAddr() string {
	return fmt.Sprintf("%s:%s", a.Host, a.GitPort)
}

func (a *App) LlmAddr() string {
	return fmt.Sprintf("%s:%s", a.Host, a.LlmPort)
}

func NewDefaultGRPCPool() *GRPCPool {
	return &GRPCPool{
		pools:   make(map[string]*grpc.ClientConn),
		LlmHost: "localhost",
		LlmPort: "00001",
		GitHost: "localhost",
		GitPort: "00002",
	}
}

type GRPCPool struct {
	mu      sync.Mutex
	pools   map[string]*grpc.ClientConn
	LlmHost string `toml:"llm_host" env:"LLM_HOST"`
	LlmPort string `toml:"llm_port" env:"LLM_PORT"`
	GitHost string `toml:"git_host" env:"GIT_HOST"`
	GitPort string `toml:"git_port" env:"GIT_PORT"`
}

// 连接池注册
func (p *GRPCPool) InitClientConn(client string) *grpc.ClientConn {
	services := map[string]string{
		"llm": fmt.Sprintf("%s:%s", p.LlmHost, p.LlmPort),
		"git": fmt.Sprintf("%s:%s", p.GitHost, p.GitPort),
	}

	if services[client] == "" {
		panic(fmt.Sprintf("service %s not found", client))
	}
	conn, err := p.GetGRPCConnection(client, services[client])
	if err != nil {
		panic(err)
	}
	return conn

}

// 获取 gRPC 连接，不用显式传入目标
func (p *GRPCPool) GetGRPCConnection(serviceName string, target string) (*grpc.ClientConn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 如果连接已存在，直接返回
	if conn, ok := p.pools[serviceName]; ok {
		return conn, nil
	}

	// 否则创建新连接
	crendital := client.NewAuthentication("admin", "123456")
	conn, err := grpc.DialContext(
		context.Background(),
		target,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(crendital),
	)
	if err != nil {
		return nil, err
	}

	// 保存到连接池
	p.pools[serviceName] = conn
	return conn, nil
}

// 获取服务连接（通过服务名）
func (p *GRPCPool) GetServiceConn(serviceName string) (*grpc.ClientConn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	conn, ok := p.pools[serviceName]
	if !ok {
		return nil, fmt.Errorf("service %s not found in connection pool", serviceName)
	}
	return conn, nil
}

// 关闭所有连接
func (p *GRPCPool) CloseAll() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for serviceName, conn := range p.pools {
		conn.Close()
		delete(p.pools, serviceName)
	}
}
