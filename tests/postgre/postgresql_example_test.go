package postgre

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"

	"go-hexagonal/config"
)

func TestMockUserDataToPostgreSQL(t *testing.T) {
	ctx := context.Background()

	var testCases = []struct {
		Name      string
		pgsqlData []string
	}{
		{
			Name: "normal test",
			pgsqlData: []string{
				"INSERT INTO users (id, email, uid) VALUES (1, 'testing@gmail.com', 'abcdefghijklmnopqrstuvwxyz12')",
			},
		},
	}

	postgresDBConf := SetupPostgreSQL(t)
	config.GlobalConfig.MigrationDir = "./migrations"

	for _, testcase := range testCases {
		t.Log("testing ", testcase.Name)

		MockPgSQLData(t, config.GlobalConfig, testcase.pgsqlData)

		connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			postgresDBConf.Username, postgresDBConf.Password, postgresDBConf.Host, postgresDBConf.Port, postgresDBConf.DbName)
		conn, err := pgx.Connect(ctx, connStr)
		if err != nil {
			t.Errorf("connect to PostgreSQL fail: %v", err)
		}
		defer conn.Close(context.TODO())

		var name, email string
		// FIXME
		err = conn.QueryRow(ctx, "SELECT name, email FROM users WHERE id = $1", 1).Scan(name, email)
		if err != nil {
			t.Errorf("query data fail: %v", err)
		}

		fmt.Println(name, email)

	}
}