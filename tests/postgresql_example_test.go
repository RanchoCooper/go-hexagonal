package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-hexagonal/config"
)

func TestMockPostgreSQLData(t *testing.T) {
	ctx := context.Background()

	var testCases = []struct {
		Name      string
		pgsqlData []string
	}{
		{
			Name: "normal test",
			pgsqlData: []string{
				"INSERT INTO users (id, name, email, uid) VALUES (1,'rancho', 'testing@gmail.com', 'abcdefghijklmnopqrstuvwxyz12')",
			},
		},
	}

	postgresDBConf := SetupPostgreSQL(t)
	config.GlobalConfig.Postgre = postgresDBConf
	config.GlobalConfig.MigrationDir = "./migrations"

	for _, testcase := range testCases {
		t.Log("testing ", testcase.Name)

		pg := MockPgSQLData(t, config.GlobalConfig, testcase.pgsqlData)

		var name, email string
		err := pg.QueryRow(ctx, "SELECT name, email FROM users WHERE id = $1", 1).Scan(&name, &email)
		if err != nil {
			t.Errorf("query data fail: %v", err)
		}
		assert.Equal(t, name, "rancho")
		assert.Equal(t, email, "testing@gmail.com")
	}
}
