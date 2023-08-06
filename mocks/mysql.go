package mocks

import (
	"context"
	"mtg-report/internal/sources/databases/mysql"

	"github.com/stretchr/testify/mock"
)

type clientMock struct {
	mock.Mock
}

func NewClientMock() *clientMock {
	return &clientMock{}
}

func (c *clientMock) ExecContext(ctx context.Context, query string, args ...interface{}) (mysql.Result, error) {
	argsMock := c.Called(ctx, query, args)
	return argsMock.Get(0).(mysql.Result), argsMock.Error(1)
}

func (c *clientMock) QueryContext(ctx context.Context, query string, args ...interface{}) (mysql.RowsScanner, error) {
	argsMock := c.Called(ctx, query, args)
	return argsMock.Get(0).(mysql.RowsScanner), argsMock.Error(1)
}

func (c *clientMock) QueryRowContext(ctx context.Context, query string, args ...interface{}) mysql.RowScanner {
	argsMock := c.Called(ctx, query, args)
	return argsMock.Get(0).(mysql.RowScanner)
}

func (c *clientMock) BeginTx(ctx context.Context, opts mysql.TxOptions) (mysql.Transaction, error) {
	argsMock := c.Called(ctx, opts)
	return argsMock.Get(0).(mysql.Transaction), argsMock.Error(1)
}

type transactionMock struct {
	mock.Mock
}

func NewTransactionMock() *transactionMock {
	return &transactionMock{}
}

func (t *transactionMock) Commit() error {
	argsMock := t.Called()
	return argsMock.Error(0)
}

func (t *transactionMock) Rollback() error {
	argsMock := t.Called()
	return argsMock.Error(0)
}

func (t *transactionMock) ExecContext(ctx context.Context, query string, args ...interface{}) (mysql.Result, error) {
	argsMock := t.Called(ctx, query, args)
	return argsMock.Get(0).(mysql.Result), argsMock.Error(1)
}

func (t *transactionMock) QueryContext(ctx context.Context, query string, args ...interface{}) (mysql.RowsScanner, error) {
	argsMock := t.Called(ctx, query, args)
	return argsMock.Get(0).(mysql.RowsScanner), argsMock.Error(1)
}

func (t *transactionMock) QueryRowContext(ctx context.Context, query string, args ...interface{}) mysql.RowScanner {
	argsMock := t.Called(ctx, query, args)
	return argsMock.Get(0).(mysql.RowScanner)
}

type rowsScannerMock struct {
	mock.Mock
}

func NewRowsScannerMock() *rowsScannerMock {
	return &rowsScannerMock{}
}

func (r *rowsScannerMock) Next() bool {
	argsMock := r.Called()
	return argsMock.Bool(0)
}

func (r *rowsScannerMock) Scan(dest ...interface{}) error {
	argsMock := r.Called(dest)
	return argsMock.Error(0)
}

func (r *rowsScannerMock) Err() error {
	argsMock := r.Called()
	return argsMock.Error(0)
}

func (r *rowsScannerMock) Close() error {
	argsMock := r.Called()
	return argsMock.Error(0)
}

type rowScannerMock struct {
	mock.Mock
}

func NewRowScannerMock() *rowScannerMock {
	return &rowScannerMock{}
}

func (r *rowScannerMock) Scan(dest ...interface{}) error {
	argsMock := r.Called()
	return argsMock.Error(0)
}

type resultMock struct {
	mock.Mock
}

func NewResultMock() *resultMock {
	return &resultMock{}
}

func (r *resultMock) LastInsertId() (int64, error) {
	argsMock := r.Called()
	return argsMock.Get(0).(int64), argsMock.Error(1)
}

func (r *resultMock) RowsAffected() (int64, error) {
	argsMock := r.Called()
	return argsMock.Get(0).(int64), argsMock.Error(1)
}

type txOptionsMock struct {
	mock.Mock
}

func NewTxOptionsMock() *txOptionsMock {
	return &txOptionsMock{}
}

func (t *txOptionsMock) IsolationLevel() mysql.Isolation {
	argsMock := t.Called()
	return argsMock.Get(0).(mysql.Isolation)
}

func (t *txOptionsMock) ReadOnly() bool {
	argsMock := t.Called()
	return argsMock.Bool(0)
}

type isolationMock struct {
	mock.Mock
}

func NewIsolationMock() *isolationMock {
	return &isolationMock{}
}

func (i *isolationMock) String() string {
	argsMock := i.Called()
	return argsMock.String(0)
}
