package cllct

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/ryym/goq/goql"
)

func InitCollectors(collectors []Collector, initConf *initConf) ([]Collector, error) {
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

	for i, cl := range clls {
		err := cl.AfterInit(initConf)
		if err != nil {
			return nil, errors.Wrapf(
				err, "failed to initialize collectors[%d] (%s)",
				i, reflect.TypeOf(cl).Elem().Name(),
			)
		}
	}
	return clls, nil
}

func fillUntakenCols(ptrs []interface{}, conf *initConf) {
	// Rows.Scan stops scanning when it encounters a nil pointer
	// in the given pointers and all subsequent pointers are ignored.
	// We need to pass a dummy pointer to prevent this.
	dummyPtr := new(interface{})
	for i, _ := range conf.ColNames {
		if conf.canTake(i) {
			ptrs[i] = dummyPtr
		}
	}
}

func ApplyCollectors(rows *sql.Rows, ptrs []interface{}, clls []Collector) error {
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

type Queryable interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type Runner struct {
	ctx   context.Context
	db    Queryable
	query goql.QueryExpr
}

func NewRunner(ctx context.Context, db Queryable, query goql.QueryExpr) *Runner {
	return &Runner{ctx, db, query}
}

func (r *Runner) Rows() (*sql.Rows, error) {
	q, err := r.query.Construct()
	if err != nil {
		return nil, err
	}
	return r.db.QueryContext(r.ctx, q.Query(), q.Args()...)
}

func (r *Runner) First(collectors ...SingleCollector) error {
	clls := make([]Collector, 0, len(collectors))
	for _, c := range collectors {
		clls = append(clls, c)
	}

	// Use WithLimits instead of Limit to avoid mutating the given query.
	return r.collect(r.query.WithLimits(1, 0), clls...)
}

func (r *Runner) Collect(collectors ...ListCollector) error {
	clls := make([]Collector, 0, len(collectors))
	for _, c := range collectors {
		clls = append(clls, c)
	}
	return r.collect(r.query, clls...)
}

func (r *Runner) collect(query goql.QueryExpr, collectors ...Collector) error {
	q, err := query.Construct()
	if err != nil {
		return err
	}

	rows, err := r.db.QueryContext(r.ctx, q.Query(), q.Args()...)
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

	ptrs := make([]interface{}, len(colNames))
	fillUntakenCols(ptrs, initConf)

	return ApplyCollectors(rows, ptrs, clls)
}

// ExecCollectorsForTest executes given collectors for
// given rows and selects.
// This is used for internal tests and not intended to
// be used for other purposes.
func ExecCollectorsForTest(
	cllcts []Collector,
	rows [][]interface{},
	selects []goql.Selection,
	colNames []string,
) error {
	if selects == nil {
		selects = make([]goql.Selection, len(colNames))
	} else {
		colNames = make([]string, len(selects))
	}

	initConf := NewInitConf(selects, colNames)
	cllcts, err := InitCollectors(cllcts, initConf)
	if err != nil {
		return err
	}

	for _, row := range rows {
		ptrs := make([]interface{}, len(selects))
		for _, cl := range cllcts {
			cl.Next(ptrs)
		}
		for i, p := range ptrs {
			if p != nil {
				reflect.ValueOf(p).Elem().Set(reflect.ValueOf(row[i]))
			}
		}
		for _, cl := range cllcts {
			cl.AfterScan(ptrs)
		}
	}

	return nil
}
