package mysql

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"go-hexagonal/internal/adapter/repository"
)

/**
 * @author Rancho
 * @date 2022/12/30
 */

type TransactionImpl struct {
}

func NewTransaction(ctx context.Context, opt *sql.TxOptions) *repository.Transaction {
	session := repository.Clients.MySQL.GetDB(ctx)
	if opt != nil {
		session = session.Begin(opt)
	}

	return &repository.Transaction{
		Session: session,
		TxOpt:   opt,
	}
}

func (t TransactionImpl) ConnDB(ctx context.Context, tr *repository.Transaction) (db *gorm.DB, err error) {
	if tr == nil {
		tr = &repository.Transaction{}
	}

	tr.Session = repository.Clients.MySQL.GetDB(ctx).Begin(tr.TxOpt)
	return tr.Session, err
}

func (t TransactionImpl) Begin(ctx context.Context, tr *repository.Transaction) {
	if tr == nil {
		tr = &repository.Transaction{}
	}
	if tr.Session == nil {
		tr.Session = repository.Clients.MySQL.GetDB(ctx).Begin(tr.TxOpt)
	}

}

func (t TransactionImpl) Commit(tr *repository.Transaction) error {
	if tr == nil {
		return errors.New("Commit with nil tr")
	}
	if tr.Session == nil {
		return errors.New("Commit with nil tr.Session")
	}

	return tr.Session.Commit().Error
}

func (t TransactionImpl) Rollback(tr *repository.Transaction) error {
	if tr == nil {
		return errors.New("Rollback with nil tr")
	}
	if tr.Session == nil {
		return errors.New("Rollback with nil tr.Session")
	}

	return tr.Session.Rollback().Error
}

var _ repository.ITransaction = &TransactionImpl{}
