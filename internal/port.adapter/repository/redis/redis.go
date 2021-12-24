package redis

import (
    "context"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/go-redis/redismock/v8"
    "github.com/pkg/errors"

    "go-hexagonal/config"
    "go-hexagonal/util/logger"
)

/**
 * @author Rancho
 * @date 2021/12/21
 */

type IRedis interface {
    GetClient() *redis.Client
    Close(ctx context.Context)
    // MockClient use only for unit test help to do  unit test without redis server
    MockClient() redismock.ClusterClientMock
}

type client struct {
    client *redis.Client
}

func (c *client) GetClient() *redis.Client {
    return c.client
}

func (c *client) Close(ctx context.Context) {
    err := c.client.Close()
    if err != nil {
        logger.Log.Errorf(ctx, "close redis client fail. err: %v", errors.WithStack(err))
    }
    logger.Log.Info(ctx, "redis client closed")
}

func (c *client) MockClient() redismock.ClusterClientMock {
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

func NewRedisClient() IRedis {
    return &client{client: NewRedis()}
}
