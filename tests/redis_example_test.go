package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockRedisData(t *testing.T) {
	var testCases = []struct {
		Name     string
		MockData map[string]string
	}{
		{
			Name: "normal test",
			MockData: map[string]string{
				"token": "mock-token",
				"name":  "rancho",
			},
		},
	}

	_, miniRedis := SetupRedis(t)

	for _, testcase := range testCases {
		t.Log("testing ", testcase.Name)

		MockRedisData(t, miniRedis, testcase.MockData)

		token, err := miniRedis.Get("token")
		assert.NoError(t, err)
		assert.Equal(t, token, "mock-token")
		name, err := miniRedis.Get("name")
		assert.NoError(t, err)
		assert.Equal(t, name, "rancho")
	}
}
