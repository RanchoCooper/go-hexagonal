package mysql

import (
    "context"
    "testing"

    "go-hexagonal/internal/domain/entity"
)

/**
 * @author Rancho
 * @date 2022/1/8
 */

var ctx = context.TODO()

func TestMain(m *testing.M) {
    db := NewExample(NewMySQLClient()).GetDB(ctx)
    _ = db.AutoMigrate(&entity.Example{})
    m.Run()
}
