package redis

import (
	"fmt"
	"go_web_bbs/settings"

	"github.com/go-redis/redis"
)

// 声明一个全局的rdb变量
var (
	client *redis.Client
	Nil    = redis.Nil
)

// 初始化连接
func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port),
		Password:     cfg.Password, // no password set
		DB:           cfg.DB,       // use default DB
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})
	_, err = client.Ping().Result()
	return
}

func Close() {
	_ = client.Close()
}
