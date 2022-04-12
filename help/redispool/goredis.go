package redispool

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/webchen/gotools/base/conf"
	"github.com/webchen/gotools/help/logs"
	"github.com/webchen/gotools/help/str"

	"github.com/go-redis/redis/v8"
)

// Ctx redis的CTX
var Ctx = context.Background()

var clientList = make(map[string]*redis.Client)

func init() {

	redisList := conf.GetConfig("redis", nil).(map[string]interface{})
	if len(redisList) == 0 {
		logs.Warning("redis config is nil", nil, false)
		return
	}

	for k, v := range redisList {
		vv := v.(map[string]interface{})

		host := vv["host"].(string)         // conf.GetConfig("redis."+k+".host", "").(string)
		port := str.Convert2U32(vv["port"]) //conf.GetConfig("redis."+k+".port", "").(string)
		db := str.Convert2U32(vv["db"])     // conf.GetConfig("redis."+k+".db", "0").(string)
		auth := vv["auth"].(string)         // conf.GetConfig("redis."+k+".auth", "").(string)
		poolSize := str.Convert2U32(vv["PoolSize"])
		minIdleConns := str.Convert2U32(vv["MinIdleConns"])
		c := redis.NewClient(&redis.Options{
			Addr:         host + ":" + strconv.FormatUint(uint64(port), 10),
			Password:     auth,    // no password set
			DB:           int(db), // use default DB
			PoolSize:     int(poolSize),
			MinIdleConns: int(minIdleConns),
			PoolTimeout:  time.Second * 2,
			IdleTimeout:  time.Second * 2,
		})
		clientList[k] = c
	}
}

// GetClient 获取对象
func GetClient(key string) *redis.Client {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil
	}
	if v, ok := clientList[key]; ok {
		return v
	}
	return nil
}
