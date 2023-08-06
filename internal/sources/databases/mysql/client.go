package mysql

import (
	"context"
	"database/sql"
	"errors"
)

type mysql struct {
	db *sql.DB
}

func New(db *sql.DB) *mysql {
	return &mysql{
		db: db,
	}
}

type transaction struct {
	tx *sql.Tx
}

func (t *transaction) Commit() error {
	return t.tx.Commit()
}

func (t *transaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *transaction) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *transaction) QueryContext(ctx context.Context, query string, args ...interface{}) (RowsScanner, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

func (t *transaction) QueryRowContext(ctx context.Context, query string, args ...interface{}) RowScanner {
	return t.tx.QueryRowContext(ctx, query, args...)
}

func (m *mysql) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error) {
	return m.db.ExecContext(ctx, query, args...)
}

func (m *mysql) QueryContext(ctx context.Context, query string, args ...interface{}) (RowsScanner, error) {
	return m.db.QueryContext(ctx, query, args...)
}

func (m *mysql) QueryRowContext(ctx context.Context, query string, args ...interface{}) RowScanner {
	return m.db.QueryRowContext(ctx, query, args...)
}

func (m *mysql) BeginTx(ctx context.Context, opts TxOptions) (Transaction, error) {
	var txOpts *sql.TxOptions
	if opts != nil {
		isolationLeval, ok := opts.IsolationLevel().(sql.IsolationLevel)
		if !ok {
			return nil, errors.New("invalid isolation level")
		}

		txOpts = &sql.TxOptions{
			Isolation: isolationLeval,
			ReadOnly:  opts.ReadOnly(),
		}
	}

	tx, err := m.db.BeginTx(ctx, txOpts)
	if err != nil {
		return nil, err
	}

	return &transaction{
		tx: tx,
	}, nil
}
