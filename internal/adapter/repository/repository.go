package repository

import (
	"context"

	"go-hexagonal/config"
	"go-hexagonal/internal/adapter/repository/mysql"
	"go-hexagonal/util/log"
)

var (
	Clients = &clients{
		MySQL: mysql.Client,
	}
	// HealthCheck *redis.HealthCheck
)

type clients struct {
	MySQL *mysql.MySQL
	// repo.IExampleRepo
	// Redis *redis.IRedis
}

func (c *clients) close(ctx context.Context) {
	if c.MySQL != nil {
		c.MySQL.Close(ctx)
	}
	// if c.Redis != nil {
	// 	c.Redis.Close(ctx)
	// }
}

type Option func(*clients)

func WithMySQL() Option {
	return func(c *clients) {
		if c.MySQL == nil {
			if config.Config.MySQL != nil {
				mysql.Client = mysql.NewMySQLClient()
			} else {
				panic("init repository fail, MySQL config is empty")
			}
		}
	}
}

func WithRedis() Option {
	return func(c *clients) {
		// 	if c.Redis == nil {
		// 		if config.Config.Redis != nil {
		// 			c.Redis = redis.NewRedisClient()
		// 		} else {
		// 			panic("init repository fail, Redis config is empty")
		// 		}
		// 	}
		// 	if HealthCheck == nil {
		// 		HealthCheck = redis.NewHealthCheck(Clients.Redis)
		// 	}
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
