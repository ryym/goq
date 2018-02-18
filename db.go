package goq

import (
	"context"
	"database/sql"

	"github.com/ryym/goq/cllct"
	"github.com/ryym/goq/dialect"
	"github.com/ryym/goq/gql"
)

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

type DB struct {
	*sql.DB
	dialect dialect.Dialect
}

func (d *DB) Close() error {
	return d.DB.Close()
}

func (d *DB) Dialect() dialect.Dialect {
	return d.dialect
}

func (d *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := d.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

func (d *DB) Query(query gql.QueryExpr) *cllct.Runner {
	return cllct.NewRunner(d.DB, query)
}

func (d *DB) Exec(query gql.QueryRoot) (sql.Result, error) {
	q := query.Construct()
	if err := q.Err(); err != nil {
		return nil, err
	}
	return d.DB.Exec(q.Query(), q.Args()...)
}

type Tx struct {
	Tx *sql.Tx
}

func (tx *Tx) Query(query gql.QueryExpr) *cllct.Runner {
	return cllct.NewRunner(tx.Tx, query)
}

func (tx *Tx) Exec(query gql.QueryRoot) (sql.Result, error) {
	q := query.Construct()
	if err := q.Err(); err != nil {
		return nil, err
	}
	return tx.Tx.Exec(q.Query(), q.Args()...)
}

func (tx *Tx) Commit() error   { return tx.Tx.Commit() }
func (tx *Tx) Rollback() error { return tx.Tx.Rollback() }
