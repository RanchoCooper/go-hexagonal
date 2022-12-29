package repository

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"

    "go-hexagonal/config"
    "go-hexagonal/internal/adapter/repository/mysql"
    "go-hexagonal/internal/adapter/repository/redis"
    "go-hexagonal/util/log"
)

/**
 * @author Rancho
 * @date 2022/1/7
 */

var ctx = context.TODO()

func TestNewRepository(t *testing.T) {
    config.Init()
    log.Init()
    Init(
        WithMySQL(),
        WithRedis(),
    )
    // mysql
    mysql.NewExample(Clients.MySQL)
    assert.NotNil(t, Example)
    assert.NotNil(t, Example.GetDB(ctx))

    // redis
    redis.NewHealthCheck(Clients.Redis)
    assert.NotNil(t, HealthCheck)
    err := HealthCheck.HealthCheck(ctx)
    assert.Nil(t, err)
}
