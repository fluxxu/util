package util

import (
	"database/sql"
	"fmt"
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

func (tx Txx) RollbackWithErrorf(f string, args ...interface{}) error {
	if txErr := tx.Rollback(); txErr != nil {
		return fmt.Errorf("%s; %s", fmt.Sprintf(f, args...), txErr)
	}
	return fmt.Errorf(f, args...)
}

func (tx Txx) Rollback() error {
	tx.BeforeRollback.Invoke()
	return tx.Tx.Rollback()
}
