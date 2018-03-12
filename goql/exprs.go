package goql

import (
	"reflect"

	"github.com/pkg/errors"
)

type aliased struct {
	exp   Expr
	alias string
}

func (a *aliased) Alias() string { return a.alias }

func (a *aliased) Apply(q *Query, ctx DBContext) {
	a.exp.Apply(q, ctx)
	q.query = append(q.query, " AS ", ctx.QuoteIdent(a.alias))
}

func (a *aliased) Selection() Selection {
	sel := a.exp.Selection()
	sel.Alias = a.alias
	return sel
}

// XXX: This should be used only in ORDER BY?
type nameExpr struct {
	name string
	ops
}

func (nm *nameExpr) init() *nameExpr {
	nm.ops = ops{nm}
	return nm
}

func (nm *nameExpr) Apply(q *Query, ctx DBContext) {
	q.query = append(q.query, ctx.QuoteIdent(nm.name))
}

func (nm *nameExpr) Selection() Selection {
	return Selection{Alias: nm.name}
}

type litExpr struct {
	val interface{}
	typ string
	ops
}

func (l *litExpr) init() *litExpr {
	l.ops = ops{l}
	return l
}

// TODO: Add no placeholder version?
func (l *litExpr) Apply(q *Query, ctx DBContext) {
	q.query = append(q.query, ctx.Placeholder(l.typ, q.args))
	q.args = append(q.args, l.val)
}

func (l *litExpr) Selection() Selection { return Selection{} }

type nullExpr struct {
	ops
}

func (n *nullExpr) init() *nullExpr {
	n.ops = ops{n}
	return n
}

func (n *nullExpr) Apply(q *Query, ctx DBContext) {
	q.query = append(q.query, "NULL")
}

func (n *nullExpr) Selection() Selection { return Selection{} }

type predExpr struct {
	AnonExpr
}

func (p *predExpr) ImplPredExpr() {}

// RawExpr is a raw string expression.
type RawExpr struct {
	sql string
	ops
}

// Allow to use it as a predicate.
func (r *RawExpr) ImplPredExpr() {}

func (r *RawExpr) init() *RawExpr {
	r.ops = ops{r}
	return r
}

func (r *RawExpr) Apply(q *Query, ctx DBContext) {
	q.query = append(q.query, r.sql)
}

func (r *RawExpr) Selection() Selection { return Selection{} }

type parensExpr struct {
	exp Expr
	ops
}

func (p *parensExpr) init() *parensExpr {
	p.ops = ops{p}
	return p
}

func (p *parensExpr) Apply(q *Query, ctx DBContext) {
	q.query = append(q.query, "(")
	p.exp.Apply(q, ctx)
	q.query = append(q.query, ")")
}

func (p *parensExpr) Selection() Selection { return p.exp.Selection() }

type inExpr struct {
	val  Expr
	exps []Expr
	not  bool
	err  error
	ops
}

func (ie *inExpr) init(in interface{}) *inExpr {
	refl := reflect.ValueOf(in)
	if refl.Type().Kind() == reflect.Slice {
		ie.exps = make([]Expr, 0, refl.Len())
		for i := 0; i < refl.Len(); i++ {
			val := refl.Index(i).Interface()
			ie.exps = append(ie.exps, lift(val))
		}
	} else if qe, ok := in.(QueryExpr); ok {
		ie.exps = []Expr{qe}
	} else {
		ie.err = errors.New("In/NotIn only accept a slice or a query")
	}

	ie.ops = ops{ie}
	return ie
}

func (ie *inExpr) Apply(q *Query, ctx DBContext) {
	if ie.err != nil {
		q.errs = append(q.errs, ie.err)
		return
	}

	ie.val.Apply(q, ctx)
	if ie.not {
		q.query = append(q.query, " NOT")
	}
	q.query = append(q.query, " IN (")
	if len(ie.exps) > 0 {
		ie.exps[0].Apply(q, ctx)
		for i := 1; i < len(ie.exps); i++ {
			q.query = append(q.query, ", ")
			ie.exps[i].Apply(q, ctx)
		}
	}
	q.query = append(q.query, ")")
}

func (ie *inExpr) Selection() Selection { return Selection{} }

type funcExpr struct {
	name string
	args []Expr
	ops
}

func (f *funcExpr) init() *funcExpr {
	f.ops = ops{f}
	return f
}

func (f *funcExpr) Apply(q *Query, ctx DBContext) {
	q.query = append(q.query, f.name+"(")
	lastIdx := len(f.args) - 1
	for i, a := range f.args {
		a.Apply(q, ctx)
		if i < lastIdx {
			q.query = append(q.query, ", ")
		}
	}
	q.query = append(q.query, ")")
}

func (f *funcExpr) Selection() Selection { return Selection{} }

// ColumnList is a special expression that holds multiple columns.
// You can pass this to Builder.Select method.
// All columns held by this struct will be selected.
type ColumnList struct {
	cols []*Column
}

func NewColumnList(cols []*Column) *ColumnList {
	return &ColumnList{cols}
}

// Columns returns the columns as a slice.
func (el *ColumnList) Columns() []*Column {
	return el.cols
}

func (el *ColumnList) Apply(q *Query, ctx DBContext) {
	if len(el.cols) == 0 {
		return
	}
	el.cols[0].Apply(q, ctx)
	for i := 1; i < len(el.cols); i++ {
		q.query = append(q.query, ", ")
		el.cols[i].Apply(q, ctx)
	}
}

func (el *ColumnList) Selection() Selection {
	panic("[INVALID] ColumnList.Selection is called")
}

type existsExpr struct {
	query QueryExpr
	not   bool
	ops
}

func (e *existsExpr) init() *existsExpr {
	e.ops = ops{e}
	return e
}

func (e *existsExpr) Apply(q *Query, ctx DBContext) {
	if e.not {
		q.query = append(q.query, "NOT ")
	}
	q.query = append(q.query, "EXISTS (")
	e.query.Apply(q, ctx)
	q.query = append(q.query, ")")
}

func (e *existsExpr) Selection() Selection { return Selection{} }
