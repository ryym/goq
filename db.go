package goq

import (
	"context"
	"database/sql"

	"github.com/ryym/goq/dialect"
)

// Open opens a database.
// The arguments are passed to the Open method of *sql.DB.
// See https://golang.org/pkg/database/sql/#Open for details.
func Open(driver, source string) (*DB, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}
	dl := dialect.New(driver)

	if dl == nil {
		dl = dialect.Generic()
	}

	return &DB{db, dl}, nil
}

// DB is a database handle which wraps *sql.DB.
// You can use Goq's query to access a DB instead of raw string SQL.
type DB struct {
	DB      *sql.DB
	dialect dialect.Dialect
}

func (d *DB) Close() error {
	return d.DB.Close()
}

// Dialect returns a dialect.Dialect for this DB.
// If this DB is not supported primarily, the dialect
// will be the generic one.
func (d *DB) Dialect() dialect.Dialect {
	return d.dialect
}

func (d *DB) Begin() (*Tx, error) {
	return d.BeginTx(context.Background(), nil)
}

func (d *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := d.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

func (d *DB) Query(query QueryExpr) *Runner {
	return d.QueryContext(context.Background(), query)
}

func (d *DB) QueryContext(ctx context.Context, query QueryExpr) *Runner {
	return NewRunner(ctx, d.DB, query)
}

func (d *DB) Exec(query QueryRoot) (sql.Result, error) {
	return d.ExecContext(context.Background(), query)
}

func (d *DB) ExecContext(ctx context.Context, query QueryRoot) (sql.Result, error) {
	q, err := query.Construct()
	if err != nil {
		return nil, err
	}
	return d.DB.ExecContext(ctx, q.Query(), q.Args()...)
}

// Tx is an in-progress database transaction which wraps *sql.Tx.
// You can use Goq's query to access a DB instead of raw string SQL.
type Tx struct {
	Tx *sql.Tx
}

func (tx *Tx) Query(query QueryExpr) *Runner {
	return NewRunner(context.Background(), tx.Tx, query)
}

func (tx *Tx) Exec(query QueryRoot) (sql.Result, error) {
	q, err := query.Construct()
	if err != nil {
		return nil, err
	}
	return tx.Tx.Exec(q.Query(), q.Args()...)
}

func (tx *Tx) Commit() error   { return tx.Tx.Commit() }
func (tx *Tx) Rollback() error { return tx.Tx.Rollback() }
