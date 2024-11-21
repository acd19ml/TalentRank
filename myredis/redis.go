package myredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

// Rdb 定义一个全局变量
var Rdb *redis.Client

// init 函数：在包被导入时自动执行
func init() {
	// 初始化 Redis 客户端
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 地址
		Password: "",               // Redis 密码（默认无密码）
		DB:       1,                // Redis 数据库（默认使用 0）
		PoolSize: 100,              // 连接池大小
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := Rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("failed to initialize Redis: %v", err)
	}

	log.Println("Redis client initialized successfully")
}
