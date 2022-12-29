package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
 * @author Rancho
 * @date 2022/12/29
 */

func TestExampleCache_HealthCheck(t *testing.T) {
	cache := NewExampleCache()
	mock := Client.MockClient()
	mock.ClearExpect()
	err := cache.HealthCheck(ctx)
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
