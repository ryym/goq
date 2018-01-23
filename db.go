package goq

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
	dl := dialect.New(driver)

	// TODO: Use generic dialect.
	if dl == nil {
		return nil, fmt.Errorf("No dialect found for %s", driver)
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

func (d *DB) Query(query gql.QueryExpr) *Collectable {
	return &Collectable{d.DB, query}
}

type Collectable struct {
	db    *sql.DB
	query gql.QueryExpr
}

func (cl *Collectable) Rows() (*sql.Rows, error) {
	q := cl.query.Construct()

	// TODO: remove debug code.
	fmt.Printf("[LOG] %s %v\n", q.Query(), q.Args())
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

	// TODO: remove debug code.
	fmt.Printf("[LOG] %s %v\n", q.Query(), q.Args())
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
	for _, cl := range collectors {
		if cl.Init(selects, colNames) {
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
