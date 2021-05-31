package redispool

import (
	"context"
	"time"

	"github.com/webchen/gotools/help/str"
	"github.com/webchen/gotools/help/tool/conf"

	"github.com/go-redis/redis/v8"
)

// Ctx redis的CTX
var Ctx = context.Background()

// Client redis对象
var Client *redis.Client

func init() {
	host := conf.GetConfig("redis.main.host", "").(string)
	port := conf.GetConfig("redis.main.port", "").(string)
	db := conf.GetConfig("redis.main.db", "0").(string)
	auth := conf.GetConfig("redis.main.auth", "").(string)
	Client = redis.NewClient(&redis.Options{
		Addr:         host + ":" + port,
		Password:     auth,                    // no password set
		DB:           int(str.String2Int(db)), // use default DB
		PoolSize:     10000,
		MinIdleConns: 1000,
		PoolTimeout:  time.Second * 2,
		IdleTimeout:  time.Second * 2,
	})
}
