package mysql

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	driver "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go-hexagonal/internal/adapter/repository"
)

/**
 * @author Rancho
 * @date 2022/12/30
 */

type TransactionImpl struct {
	*repository.Transaction
}

func (t TransactionImpl) Begin(ctx context.Context, tr *repository.Transaction) (*repository.Transaction, error) {
	gormDB := repository.Clients.MySQL.GetDB(ctx)
	if tr == nil {
		sqlDB, err := gormDB.DB()
		if err != nil {
			return nil, errors.Wrap(err, "begin transaction fail when sqlDB")
		}
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			return nil, errors.Wrap(err, "begin transaction fail when beginTx")
		}
		tr = &repository.Transaction{
			Tx:     tx,
			TxOpts: nil,
		}
	}
	if tr.TxOpts != nil {
		// begin new tx with opts
		gormDB = gormDB.Begin(tr.TxOpts...)
		tx, ok := gormDB.Statement.ConnPool.(*sql.Tx)
		if !ok {
			return nil, errors.New("Begin fail when convert gorm.ConnPool to sql.Tx")
		}
		tr.Tx = tx
	}
	return tr, nil
}

func (t TransactionImpl) Commit(tr *repository.Transaction) error {
	if tr == nil {
		return errors.New("Commit with nil Transaction")
	}
	if tr.Tx == nil {
		return errors.New("Commit with nil Transaction.Tx")
	}
	return tr.Tx.Commit()
}

func (t TransactionImpl) RollBack(tr *repository.Transaction) error {
	if tr == nil {
		return errors.New("RollBack with nil Transaction")
	}
	if tr.Tx == nil {
		return errors.New("RollBack with nil Transaction.Tx")
	}
	return tr.Tx.Rollback()
}

func (t TransactionImpl) ConnDB(ctx context.Context, tr *repository.Transaction) (db *gorm.DB, err error) {
	// prepare transaction
	if tr == nil {
		tr, err = t.Begin(ctx, nil)
		if err != nil {
			return nil, errors.Wrap(err, "transaction begin fail")
		}
		t.Transaction = tr
	}

	// open db conn
	db, err = gorm.Open(driver.New(driver.Config{Conn: tr.Tx}))
	if err != nil {
		return nil, errors.Wrap(err, "open DB conn fail")
	}

	return db, err
}

func (t TransactionImpl) AfterCreate(tx *gorm.DB) error {
	// return t.Commit(t.Transaction)
	return tx.Commit().Error
}

var _ repository.ITransaction = &TransactionImpl{}
