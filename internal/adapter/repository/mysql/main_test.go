package mysql

import (
    "context"
    "testing"

    "go-hexagonal/config"
    "go-hexagonal/internal/domain/entity"
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

    db := NewExample(NewMySQLClient()).GetDB(ctx)
    _ = db.AutoMigrate(&entity.Example{})
    m.Run()
}
