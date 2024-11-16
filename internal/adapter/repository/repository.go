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
}

type StoreType string

const (
	MySQLStore      StoreType = "MySQL"
	RedisStore      StoreType = "Redis"
	MongoStore      StoreType = "Mongo"
	PostgreSQLStore StoreType = "PostgreSQL"
)

type clients struct {
	MySQL *MySQL
	Redis *Redis
}

type Option func(*clients)

func (tr *Transaction) Conn(ctx context.Context) *gorm.DB {
	if tr == nil {
		// init transaction with default session
		return Clients.MySQL.GetDB(ctx)
	}
	if tr.Session == nil {
		// begin new with TxOpt
		tr.Session = Clients.MySQL.GetDB(ctx).Begin(tr.TxOpt)
	}

	return tr.Session
}

func NewTransaction(ctx context.Context, store StoreType, opt *sql.TxOptions) *Transaction {
	tr := &Transaction{TxOpt: opt}

	if store == MySQLStore {
		session := Clients.MySQL.GetDB(ctx)
		if opt != nil {
			session = session.Begin(opt)
		}
		tr.Session = session
	} else if store == RedisStore {
		// TODO
	} else if store == MongoStore {
		// TODO
	} else if store == PostgreSQLStore {
		// TODO
	}

	return tr
}

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
			if config.GlobalConfig.MySQL == nil {
				panic("repository init fail, MySQL config is empty")
			}
			c.MySQL = NewMySQLClient()
		}
	}
}

func WithRedis() Option {
	return func(c *clients) {
		if c.Redis == nil {
			if config.GlobalConfig.Redis == nil {
				panic("repository init fail, Redis config is empty")
			}
			c.Redis = NewRedisClient()
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
