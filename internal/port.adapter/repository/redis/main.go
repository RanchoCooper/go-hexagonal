package redis

import (
    "sync"
    "time"

    "github.com/go-redis/redis/v8"

    "go-hexagonal/config"
)

/**
 * @author Rancho
 * @date 2021/12/21
 */

var (
    once  sync.Once
    Redis *RedisRepository
)

type RedisRepository struct {
    client *redis.Client
}

func init() {
    once.Do(func() {
        Redis = NewRedisRepository()
    })
}

func NewRedisDB() *redis.Client {
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

func NewRedisRepository() *RedisRepository {
    return &RedisRepository{client: NewRedisDB()}
}
