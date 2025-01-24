package redis

import (
	"context"
	"errors"

	"go-hexagonal/adapter/repository"
	"go-hexagonal/domain/repo"
)

func NewExampleCache() *ExampleCache {
	return &ExampleCache{}
}

type ExampleCache struct {
}

var _ repo.IExampleCacheRepo = &ExampleCache{}

func (h ExampleCache) HealthCheck(ctx context.Context) error {
	pong := repository.Clients.Redis.GetClient().Ping(ctx).String()
	if pong != "ping: PONG" {
		return errors.New("ping redis got invalid response: " + pong)
	}
	return nil
}
