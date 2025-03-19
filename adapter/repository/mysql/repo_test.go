package mysql

import (
	"context"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestTransaction_Conn(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	r := NewMysqlRepo(db)

	tx := r.Conn(context.Background())
	if tx.Error != nil {
		t.Errorf("Expected no error, got %v", tx.Error)
	}

	if err := r.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

type MysqlRepo struct {
	db *gorm.DB
}

func NewMysqlRepo(db *gorm.DB) *MysqlRepo {
	return &MysqlRepo{db: db}
}

func (r *MysqlRepo) Conn(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *MysqlRepo) Close() error {
	sqlDB, _ := r.db.DB()
	return sqlDB.Close()
}
