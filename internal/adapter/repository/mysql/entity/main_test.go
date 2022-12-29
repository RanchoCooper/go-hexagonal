package entity

import (
	"context"
	"testing"

	"go-hexagonal/config"
	"go-hexagonal/internal/adapter/repository/mysql"
	"go-hexagonal/util/log"
)

/**
 * @author Rancho
 * @date 2022/1/8
 */

var ctx = context.TODO()

func TestMain(m *testing.M) {
	config.Init()
	log.Init()

	mysql.Client = mysql.NewMySQLClient()
	_ = mysql.Client.GetDB(ctx).AutoMigrate(&Example{})
	m.Run()
}
