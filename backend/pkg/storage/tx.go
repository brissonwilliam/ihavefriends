package storage

import (
	"errors"
	"github.com/jmoiron/sqlx"
)

// TxProvider is a very generic abstraction of sql transactions
// It's mostly to simplify usage of transactions and storage in usecases
type TxProvider interface {
	Begin() (UnitOfWork, error)
	Close(uow UnitOfWork, err *error)
}

func NewTxProvider(db *sqlx.DB) TxProvider {
	return defaultTxProvider{
		db: db,
	}
}

type defaultTxProvider struct {
	db *sqlx.DB
}

func (tp defaultTxProvider) Begin() (UnitOfWork, error) {
	return tp.db.Beginx()
}

// Close is a helper that automatically rolls back or commits
// the unit of work, whether there is a panic, error, or everything is alright.
func (tp defaultTxProvider) Close(uow UnitOfWork, err *error) {
	p := recover()

	if p != nil {
		if e, ok := p.(error); ok {
			// override error value to be of panic error
			*err = e
		} else {
			*err = errors.New("Rolling back due to unexpected panic")
		}
	}

	if *err != nil {
		_ = uow.Rollback()
	} else {
		_ = uow.Commit()
	}

}

// UnitOfWork is a very generic way to wrap transactions
// and avoid directly managing them at use-case level.
type UnitOfWork interface {
	Rollback() error
	Commit() error
}

func UnitAsTransaction(uow UnitOfWork) *sqlx.Tx {
	return uow.(*sqlx.Tx)
}
