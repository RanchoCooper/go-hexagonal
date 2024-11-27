package tests

import (
	"strconv"
	"testing"

	"github.com/alicebob/miniredis/v2"

	"go-hexagonal/config"
)

func SetupRedis(t *testing.T) (redisConf *config.RedisConfig, s *miniredis.Miniredis) {
	s = miniredis.RunT(t)

	redisPort, err := strconv.Atoi(s.Port())
	if err != nil {
		t.Fatalf("failed to get redis post, err: %s", err)
	}

	return &config.RedisConfig{
		Host: s.Host(),
		Port: redisPort,
	}, s
}

func MockRedisData(t *testing.T, miniRedis *miniredis.Miniredis, data map[string]string) {
	miniRedis.FlushAll()

	for k, v := range data {
		err := miniRedis.Set(k, v)
		if err != nil {
			t.Fatalf("mock redis data fail, k: %s, v: %s, err: %s", k, v, err)
		}
	}
}
