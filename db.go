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
