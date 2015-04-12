package util

import (
	"database/sql"
	"github.com/lann/squirrel"
)

type SquirrelDBProxy struct {
	*sql.DB
}

func (p *SquirrelDBProxy) QueryRow(sql string, args ...interface{}) squirrel.RowScanner {
	return p.DB.QueryRow(sql, args...)
}

type Txx struct {
	*sql.Tx
	OnCommit   Delegate
	OnRollback Delegate
}

func (tx Txx) Commit() error {
	tx.OnCommit.Invoke()
	return tx.Tx.Commit()
}

func (tx Txx) Rollback() error {
	tx.OnRollback.Invoke()
	return tx.Tx.Rollback()
}
