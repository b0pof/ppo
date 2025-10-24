package user

import (
	"context"
	"database/sql"
)

type database interface {
	ExecContext(ctx context.Context, q string, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, q string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, q string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}
