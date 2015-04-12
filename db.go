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
	AfterCommit    Delegate
	BeforeRollback Delegate
}

func (tx Txx) Commit() error {
	if err := tx.Tx.Commit(); err != nil {
		return err
	}
	tx.AfterCommit.Invoke()
	return nil
}

func (tx Txx) Rollback() error {
	tx.BeforeRollback.Invoke()
	return tx.Tx.Rollback()
}
