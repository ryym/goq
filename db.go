package main

import (
	"database/sql"
	"fmt"

	"github.com/ryym/goq/cllct"
	"github.com/ryym/goq/dialect"
	"github.com/ryym/goq/gql"
)

func Open(driver, source string) (*DB, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}
	dl, err := dialect.New(driver)

	// TODO: Use generic dialect.
	if err != nil {
		return nil, err
	}

	return &DB{db, dl}, nil
}

type DB struct {
	*sql.DB
	dialect dialect.Dialect
}

func (d *DB) QueryBuilder() *Builder {
	return &Builder{
		Builder:        gql.NewBuilder(d.dialect),
		CollectorMaker: cllct.NewMaker(),
	}
}

func (d *DB) Query(query gql.QueryExpr) *Collectable {
	q := query.Construct()
	return &Collectable{d.DB, q.Query(), q.Args()}
}

type Collectable struct {
	db    *sql.DB
	query string
	args  []interface{}
}

func (cl *Collectable) Rows() (*sql.Rows, error) {
	// TODO: remove debug code.
	fmt.Printf("[LOG] %s %v\n", cl.query, cl.args)
	return cl.db.Query(cl.query, cl.args...)
}
