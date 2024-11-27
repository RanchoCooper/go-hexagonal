package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-hexagonal/config"
)

func TestMockMySQLData(t *testing.T) {
	var testCases = []struct {
		Name    string
		sqlData []string
	}{
		{
			Name: "normal test",
			sqlData: []string{
				"INSERT INTO users (id, name, email, uid) VALUES (1,'rancho', 'testing@gmail.com', 'abcdefghijklmnopqrstuvwxyz12')",
			},
		},
	}

	mysqlDBConf := SetupMySQL(t)
	config.GlobalConfig.MySQL = mysqlDBConf
	config.GlobalConfig.MigrationDir = "./migrations"

	for _, testcase := range testCases {
		t.Log("testing ", testcase.Name)

		db := MockMySQLData(t, config.GlobalConfig, testcase.sqlData)

		type UserVO struct {
			Id    int
			Name  string
			Email string
		}
		user := UserVO{}
		tx := db.Raw("SELECT name, email from users where id = ?", 1).Scan(&user)
		if tx.Error != nil {
			t.Errorf("query data fail: %v", tx.Error)
		}
		assert.Equal(t, "rancho", user.Name)
		assert.Equal(t, "testing@gmail.com", user.Email)
	}
}
