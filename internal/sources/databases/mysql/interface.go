package mysql

import (
	"context"
)

type RowScanner interface {
	Scan(dest ...interface{}) error
}

type RowsScanner interface {
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
	Close() error
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type Transaction interface {
	Commit() error
	Rollback() error
	ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (RowsScanner, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) RowScanner
}

type TxOptions interface {
	IsolationLevel() Isolation
	ReadOnly() bool
}

type Isolation interface {
	String() string
}

type Client interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (RowsScanner, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) RowScanner
	BeginTx(ctx context.Context, opts TxOptions) (Transaction, error)
}
