package redis

import (
	"context"
	"testing"

	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

/**
 * @author Rancho
 * @date 2022/12/29
 */

var ctx = context.TODO()

func TestMain(m *testing.M) {
	config.Init()
	log.Init()

	Client = NewRedisClient()
	m.Run()
}
