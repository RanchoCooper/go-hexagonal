package repository

import (
	"context"
	"testing"

	"go-hexagonal/config"
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

	Init(WithMySQL(), WithRedis())
	// assert.Nil(t, err)
	// assert.NotNil(t, model)

	Close(ctx)
	// redis
}
