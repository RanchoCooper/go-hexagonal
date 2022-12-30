package repository

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

var Clients = &clients{}

type Transaction struct {
	Session *gorm.DB
	TxOpt   *sql.TxOptions
	// Tx      *sql.Tx
}

type ITransaction interface {
	Begin(context.Context, *Transaction)
	Commit(*Transaction) error
	Rollback(*Transaction) error
}

type clients struct {
	MySQL *MySQL
	Redis *Redis
}

type Option func(*clients)

func (c *clients) close(ctx context.Context) {
	if c.MySQL != nil {
		c.MySQL.Close(ctx)
	}
	if c.Redis != nil {
		c.Redis.Close(ctx)
	}
}

func WithMySQL() Option {
	return func(c *clients) {
		if c.MySQL == nil {
			if config.Config.MySQL != nil {
				c.MySQL = NewMySQLClient()
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
				c.Redis = NewRedisClient()
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
