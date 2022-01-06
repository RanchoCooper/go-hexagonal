package repository

import (
    "context"

    "go-hexagonal/config"
    "go-hexagonal/internal/port.adapter/repository/mysql"
    "go-hexagonal/internal/port.adapter/repository/redis"
    "go-hexagonal/util/logger"
)

var (
    Clients     = &client{}
    Example     *mysql.Example
    HealthCheck *redis.HealthCheck
)

type client struct {
    MySQL mysql.IMySQL
    Redis redis.IRedis
}

func (c *client) close(ctx context.Context) {
    if c.MySQL != nil {
        c.MySQL.Close(ctx)
    }
    if c.Redis != nil {
        c.Redis.Close(ctx)
    }
}

type Option func(*client)

func WithMySQL(ctx context.Context) Option {
    return func(c *client) {
        if c.MySQL == nil {
            if config.Config.MySQL != nil {
                c.MySQL = mysql.NewMySQLClient()
            } else {
                panic("init repository with empty MySQL config")
            }
        }
        // inject repository
        if Example == nil {
            Example = mysql.NewExample(Clients.MySQL)
        }
    }
}

func WithRedis(ctx context.Context) Option {
    return func(c *client) {
        if c.Redis == nil {
            if config.Config.Redis != nil {
                c.Redis = redis.NewRedisClient()
            } else {
                panic("init repository with empty Redis config")
            }
        }
        if HealthCheck == nil {
            HealthCheck = redis.NewHealthCheck(Clients.Redis)
        }
    }
}

func Init(opts ...Option) {
    for _, opt := range opts {
        opt(Clients)
    }
    logger.Log.Info(context.Background(), "repository init successfully")
}

func Close(ctx context.Context) {
    Clients.close(ctx)
    logger.Log.Info(ctx, "repository is closed.")
}
