package review

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type database interface {
	ExecContext(ctx context.Context, q string, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, q string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, q string, args ...interface{}) error
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}
