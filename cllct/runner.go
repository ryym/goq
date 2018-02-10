package cllct

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/ryym/goq/gql"
)

func InitCollectors(collectors []Collector, initConf *InitConf) ([]Collector, error) {
	clls := make([]Collector, 0, len(collectors))
	for i, cl := range collectors {
		ok, err := cl.Init(initConf)
		if err != nil {
			return nil, errors.Wrapf(
				err, "failed to initialize collectors[%d] (%s)",
				i, reflect.TypeOf(cl).Elem().Name(),
			)
		}
		if ok {
			clls = append(clls, cl)
		}
	}
	return clls, nil
}

func ApplyCollectors(rows *sql.Rows, colNames []string, clls []Collector) error {
	// Rows.Scan stops scanning when it encounters a nil pointer
	// in the given pointers and all subsequent pointers are ignored.
	// We need to pass a dummy pointer to prevent this.
	dummyPtr := new(interface{})
	ptrs := make([]interface{}, len(colNames))

	for rows.Next() {
		for _, cl := range clls {
			cl.Next(ptrs)
		}

		for i := 0; i < len(ptrs); i++ {
			if ptrs[i] == nil {
				ptrs[i] = dummyPtr
			}
		}

		rows.Scan(ptrs...)

		for _, cl := range clls {
			cl.AfterScan(ptrs)
		}
	}

	return rows.Err()
}

type Queryable interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type Runner struct {
	db    Queryable
	query gql.QueryExpr
}

func NewRunner(db Queryable, query gql.QueryExpr) *Runner {
	return &Runner{db, query}
}

func (r *Runner) Rows() (*sql.Rows, error) {
	q := r.query.Construct()
	return r.db.Query(q.Query(), q.Args()...)
}

func (r *Runner) First(collectors ...SingleCollector) error {
	clls := make([]Collector, len(collectors))
	for i, c := range collectors {
		clls[i] = c
	}

	// Use WithLimits instead of Limit to avoid mutating the given query.
	return r.collect(r.query.WithLimits(1, 0), clls...)
}

func (r *Runner) Collect(collectors ...ListCollector) error {
	clls := make([]Collector, len(collectors))
	for i, c := range collectors {
		clls[i] = c
	}
	return r.collect(r.query, clls...)
}

func (r *Runner) collect(query gql.QueryExpr, collectors ...Collector) error {
	q := query.Construct()
	rows, err := r.db.Query(q.Query(), q.Args()...)
	if err != nil {
		return err
	}
	defer rows.Close()

	selects := query.Selections()
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

	initConf := NewInitConf(selects, colNames)
	clls, err := InitCollectors(collectors, initConf)
	if err != nil {
		return err
	}

	return ApplyCollectors(rows, colNames, clls)
}
