package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ctx = context.TODO()

func TestNewRepository(t *testing.T) {
	Init(WithMySQL(), WithRedis())

	Close(ctx)
}

func TestTransaction_Conn(t *testing.T) {
	Init(WithMySQL(), WithRedis())

	t.Run("nil caller", func(t *testing.T) {
		var tr *Transaction
		db := tr.Conn(ctx)
		assert.NotNil(t, db)
	})

	t.Run("with empty session", func(t *testing.T) {
		tr := NewTransaction(ctx,
			MySQLStore,
			nil,
		)
		tr.Session = nil
		db := tr.Conn(ctx)
		assert.NotNil(t, db)
	})
	t.Run("with opt", func(t *testing.T) {
		tr := NewTransaction(ctx,
			MySQLStore,
			&sql.TxOptions{
				Isolation: sql.LevelReadUncommitted,
				ReadOnly:  false,
			},
		)
		db := tr.Conn(ctx)
		assert.NotNil(t, db)
	})
}
