package initialize

import (
	"context"
	"gin-vue-admin/config"
	"gin-vue-admin/global"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

func initRedisClient(redisCfg config.Redis) (redis.UniversalClient, error) {
	var client redis.UniversalClient
	// 使用集群模式
	if redisCfg.UseCluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    redisCfg.ClusterAddrs,
			Password: redisCfg.Password,
		})
	} else {
		// 使用单例模式
		client = redis.NewClient(&redis.Options{
			Addr:     redisCfg.Addr,
			Password: redisCfg.Password,
			DB:       redisCfg.DB,
		})
	}
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		logx.Error("redis connect ping failed, err:", logx.Field("name", redisCfg.Name), logx.Field("err", err))
		return nil, err
	}

	logx.Info("redis connect ping response:", logx.Field("name", redisCfg.Name), logx.Field("pong", pong))
	return client, nil
}

func Redis() {
	redisClient, err := initRedisClient(global.GVA_CONFIG.Redis)
	if err != nil {
		panic(err)
	}
	global.GVA_REDIS = redisClient
}

func RedisList() {
	redisMap := make(map[string]redis.UniversalClient)

	for _, redisCfg := range global.GVA_CONFIG.RedisList {
		client, err := initRedisClient(redisCfg)
		if err != nil {
			panic(err)
		}
		redisMap[redisCfg.Name] = client
	}

	global.GVA_REDISList = redisMap
}
