package goq

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/ryym/goq/goql"
)

// InitCollectors initializes collectors.
// This is used internally.
func InitCollectors(collectors []Collector, initConf *initConf) ([]Collector, error) {
	clls := make([]Collector, 0, len(collectors))
	for i, cl := range collectors {
		ok, err := cl.init(initConf)
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
		err := cl.afterinit(initConf)
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

func applyCollectors(rows *sql.Rows, ptrs []interface{}, clls []Collector) (nScans int, err error) {
	nScans = 0
	for rows.Next() {
		for _, cl := range clls {
			cl.next(ptrs)
		}

		// TODO: Should not ignore scan errors.
		_ = rows.Scan(ptrs...)
		nScans += 1

		for _, cl := range clls {
			cl.afterScan(ptrs)
		}
	}

	return nScans, nil
}

// Queryable represents an interface to issue an query.
// Builtin *sql.DB implements this interface.
type Queryable interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// Runner runs a query with given collectors to collect result rows.
type Runner struct {
	ctx   context.Context
	db    Queryable
	query goql.QueryExpr
}

func NewRunner(ctx context.Context, db Queryable, query goql.QueryExpr) *Runner {
	return &Runner{ctx, db, query}
}

// Rows return *sql.Rows directly.
func (r *Runner) Rows() (*sql.Rows, error) {
	q, err := r.query.Construct()
	if err != nil {
		return nil, err
	}
	return r.db.QueryContext(r.ctx, q.Query(), q.Args()...)
}

// First executes given single collectors.
func (r *Runner) First(collectors ...SingleCollector) error {
	clls := make([]Collector, 0, len(collectors))
	for _, c := range collectors {
		clls = append(clls, c)
	}

	// Use WithLimits instead of Limit to avoid mutating the given query.
	nScans, err := r.collect(r.query.WithLimits(1, 0), clls...)
	if err != nil {
		return err
	}
	if nScans == 0 {
		return ErrNoRows
	}
	return nil
}

// First executes given list collectors.
func (r *Runner) Collect(collectors ...ListCollector) error {
	clls := make([]Collector, 0, len(collectors))
	for _, c := range collectors {
		clls = append(clls, c)
	}
	_, err := r.collect(r.query, clls...)
	return err
}

func (r *Runner) collect(query goql.QueryExpr, collectors ...Collector) (nScans int, err error) {
	q, err := query.Construct()
	if err != nil {
		return 0, err
	}

	rows, err := r.db.QueryContext(r.ctx, q.Query(), q.Args()...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	selects := query.Selections()
	colNames, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	if len(colNames) != len(selects) {
		return 0, fmt.Errorf(
			"[goq] selections mismatch: colNames: %d, selects: %d",
			len(colNames),
			len(selects),
		)
	}

	initConf := NewInitConf(selects, colNames)
	clls, err := InitCollectors(collectors, initConf)
	if err != nil {
		return 0, err
	}

	ptrs := make([]interface{}, len(colNames))
	fillUntakenCols(ptrs, initConf)

	return applyCollectors(rows, ptrs, clls)
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
			cl.next(ptrs)
		}
		for i, p := range ptrs {
			if p != nil {
				reflect.ValueOf(p).Elem().Set(reflect.ValueOf(row[i]))
			}
		}
		for _, cl := range cllcts {
			cl.afterScan(ptrs)
		}
	}

	return nil
}
