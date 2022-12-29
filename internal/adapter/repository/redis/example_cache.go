package redis

import (
	"context"
	"errors"

	"go-hexagonal/internal/domain/repo"
)

/**
 * @author Rancho
 * @date 2022/1/7
 */

func NewExampleCache() *ExampleCache {
	return &ExampleCache{}
}

type ExampleCache struct {
}

var _ repo.IExampleCacheRepo = &ExampleCache{}

func (h ExampleCache) HealthCheck(ctx context.Context) error {
	pong := Client.RedisDB.Ping(ctx).String()
	if pong != "ping: PONG" {
		return errors.New("ping redis got invalid response: " + pong)
	}
	return nil
}
