package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-hexagonal/internal/adapter/repository"
)

/**
 * @author Rancho
 * @date 2022/12/29
 */

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
