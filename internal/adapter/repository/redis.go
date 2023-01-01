package repository

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

// var (
// 	Client *Redis
// )

type Redis struct {
	db *redis.Client
}

func NewRedisClient() *Redis {
	return &Redis{db: newRedisConn()}
}

func (r *Redis) GetClient() *redis.Client {
	return r.db
}

func (r *Redis) Close(ctx context.Context) {
	if err := r.db.Close(); err != nil {
		log.SugaredLogger.Errorf("close redis client fail. err: %s", err.Error())
	}
	log.Logger.Info("redis client closed")
}

func (r *Redis) MockClient() redismock.ClientMock {
	db, mock := redismock.NewClientMock()
	r.db = db
	return mock
}

func newRedisConn() *redis.Client {
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
