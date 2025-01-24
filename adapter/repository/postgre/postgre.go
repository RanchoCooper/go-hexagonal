package postgre

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"go-hexagonal/config"
)

func NewConnPool(conf *config.PostgreSQLConfig) (*pgxpool.Pool, error) {
	pgxPoolConfig, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DbName,
	))
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolConfig)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
