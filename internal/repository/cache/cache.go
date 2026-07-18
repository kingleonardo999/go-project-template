package cache

import (
	"context"

	"go-project-template/internal/config"

	"github.com/redis/go-redis/v9"
)

// RDB 全局 Redis 客户端
var RDB *redis.Client

// Init 初始化 Redis 连接
func Init(cfg config.RedisConfig) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	// 验证连接
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic("redis connect error: " + err.Error())
	}
	RDB = rdb
}

// Close 关闭 Redis 连接
func Close() {
	if RDB != nil {
		RDB.Close()
	}
}
