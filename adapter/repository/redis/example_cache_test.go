package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-hexagonal/adapter/repository"
)

func TestExampleCache_HealthCheck(t *testing.T) {
	cache := NewExampleCache()
	mock := repository.Clients.Redis.MockClient()
	mock.ClearExpect()
	mock.ExpectPing().SetVal("PONG")

	err := cache.HealthCheck(ctx)
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
