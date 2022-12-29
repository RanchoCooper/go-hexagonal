package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"

	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

/**
 * @author Rancho
 * @date 2021/12/21
 */

var (
	Client *Redis
)

type Redis struct {
	RedisDB *redis.Client
}

func (r *Redis) GetClient() *redis.Client {
	return r.RedisDB
}

func (r *Redis) Close(ctx context.Context) {
	err := r.RedisDB.Close()
	if err != nil {
		log.SugaredLogger.Errorf("close redis client fail. err: %s", err.Error())
	}
	log.Logger.Info("redis client closed")
}

func (r *Redis) MockClient() redismock.ClusterClientMock {
	// FIXME unverified
	_, mock := redismock.NewClusterMock()
	return mock
}

func NewRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         config.Config.Redis.Addr,
		Username:     config.Config.Redis.UserName,
		Password:     config.Config.Redis.Password,
		DB:           config.Config.Redis.DB,
		PoolSize:     config.Config.Redis.PoolSize,
		MinIdleConns: config.Config.Redis.MinIdleConns,
		IdleTimeout:  time.Duration(config.Config.Redis.IdleTimeout) * time.Second,
	})
}

func NewRedisClient() *Redis {
	return &Redis{RedisDB: NewRedis()}
}
