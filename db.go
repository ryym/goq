package goq

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
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

func (d *DB) Query(query gql.QueryExpr) *Collectable {
	return &Collectable{d.DB, query}
}

type Tx struct {
	Tx *sql.Tx
}

func (tx *Tx) Query(query gql.QueryExpr) *Collectable {
	return &Collectable{tx.Tx, query}
}

func (tx *Tx) Commit() error   { return tx.Tx.Commit() }
func (tx *Tx) Rollback() error { return tx.Tx.Rollback() }

type Queryable interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type Collectable struct {
	db    Queryable
	query gql.QueryExpr
}

func (cl *Collectable) Rows() (*sql.Rows, error) {
	q := cl.query.Construct()
	return cl.db.Query(q.Query(), q.Args()...)
}

func (cl *Collectable) First(collectors ...cllct.SingleCollector) error {
	clls := make([]cllct.Collector, len(collectors))
	for i, c := range collectors {
		clls[i] = c
	}
	return cl.collect(cl.query.Limit(1), clls...)
}

func (cl *Collectable) Collect(collectors ...cllct.ListCollector) error {
	clls := make([]cllct.Collector, len(collectors))
	for i, c := range collectors {
		clls[i] = c
	}
	return cl.collect(cl.query, clls...)
}

func (cl *Collectable) collect(query gql.QueryExpr, collectors ...cllct.Collector) error {
	q := query.Construct()
	rows, err := cl.db.Query(q.Query(), q.Args()...)
	if err != nil {
		return err
	}
	defer rows.Close()

	selects := cl.query.Selections()
	colNames, err := rows.Columns()
	if err != nil {
		return err
	}

	if len(colNames) != len(selects) {
		return fmt.Errorf(
			"[goq] selections mismatch: colNames: %d, selects: %d",
			len(colNames),
			len(selects),
		)
	}

	clls := make([]cllct.Collector, 0, len(collectors))
	initConf := cllct.NewInitConf(selects, colNames)
	for i, cl := range collectors {
		ok, err := cl.Init(initConf)
		if err != nil {
			return errors.Wrapf(
				err, "failed to initialize collectors[%d] (%s)",
				i, reflect.TypeOf(cl).Name(),
			)
		}
		if ok {
			clls = append(clls, cl)
		}
	}

	ptrs := make([]interface{}, len(colNames))
	for rows.Next() {
		for _, cl := range clls {
			cl.Next(ptrs)
		}
		rows.Scan(ptrs...)
		for _, cl := range clls {
			cl.AfterScan(ptrs)
		}
	}

	return rows.Err()
}
