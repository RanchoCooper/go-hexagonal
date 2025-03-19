package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"

	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

type Redis struct {
	db *redis.Client
}

func NewRedisClient() *Redis {
	return &Redis{db: newRedisConn()}
}

func (r *Redis) GetClient() *redis.Client {
	return r.db
}

func (r *Redis) Close(ctx context.Context) error {
	if err := r.db.Close(); err != nil {
		// log.SugaredLogger.Errorf("close redis client fail. err: %s", err.Error())
		return err
	}
	log.Logger.Info("redis client closed")
	return nil
}

func (r *Redis) MockClient() redismock.ClientMock {
	db, mock := redismock.NewClientMock()
	r.db = db
	return mock
}

func newRedisConn() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         config.GlobalConfig.Redis.Host,
		Password:     config.GlobalConfig.Redis.Password,
		DB:           config.GlobalConfig.Redis.DB,
		PoolSize:     config.GlobalConfig.Redis.PoolSize,
		MinIdleConns: config.GlobalConfig.Redis.MinIdleConns,
		IdleTimeout:  time.Duration(config.GlobalConfig.Redis.IdleTimeout) * time.Second,
	})
}
