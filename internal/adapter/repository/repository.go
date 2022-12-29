package repository

import (
	"context"

	"go-hexagonal/config"
	"go-hexagonal/internal/adapter/repository/mysql"
	"go-hexagonal/internal/adapter/repository/redis"
	"go-hexagonal/util/log"
)

var Clients = &clients{}

type clients struct {
	MySQL *mysql.MySQL
	Redis *redis.Redis
}

func (c *clients) close(ctx context.Context) {
	if c.MySQL != nil {
		c.MySQL.Close(ctx)
	}
	if c.Redis != nil {
		c.Redis.Close(ctx)
	}
}

type Option func(*clients)

func WithMySQL() Option {
	return func(c *clients) {
		if c.MySQL == nil {
			if config.Config.MySQL != nil {
				mysql.Client = mysql.NewMySQLClient()
				c.MySQL = mysql.Client
			} else {
				panic("init repository fail, MySQL config is empty")
			}
		}
	}
}

func WithRedis() Option {
	return func(c *clients) {
		if c.Redis == nil {
			if config.Config.Redis != nil {
				redis.Client = redis.NewRedisClient()
				c.Redis = redis.Client
			} else {
				panic("init repository fail, Redis config is empty")
			}
		}
	}
}

func Init(opts ...Option) {
	for _, opt := range opts {
		opt(Clients)
	}
	log.Logger.Info("repository init successfully")
}

func Close(ctx context.Context) {
	Clients.close(ctx)
	log.Logger.Info("repository closed")
}
