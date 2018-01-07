package gnr

import (
	"fmt"
	"strings"

	q "github.com/ryym/goq/query"
)

type columnExpr struct {
	name       string
	tableName  string
	structName string
	fieldName  string
}

func (c *columnExpr) ToQuery() q.Query {
	return q.Query{
		fmt.Sprintf("%s.%s", c.tableName, c.name),
		[]interface{}{},
	}
}

func (c *columnExpr) ToSelectItem() q.SelectItem {
	return q.SelectItem{
		ColumnName: c.name,
		TableName:  c.tableName,
		StructName: c.structName,
		FieldName:  c.fieldName,
	}
}

type parensExpr struct {
	exp q.Expr
}

func (p *parensExpr) ToQuery() q.Query {
	qr := p.exp.ToQuery()
	return q.Query{fmt.Sprintf("(%s)", qr.Query), qr.Args}
}

func (p *parensExpr) ToSelectItem() q.SelectItem { return p.exp.ToSelectItem() }

type litExpr struct {
	val interface{}
}

func (l *litExpr) ToQuery() q.Query {
	return q.Query{"?", []interface{}{l.val}}
}

func (l *litExpr) ToSelectItem() q.SelectItem { return q.SelectItem{} }

type exprAliased struct {
	expr  q.Expr
	alias string
}

func (e *exprAliased) Alias() string { return e.alias }
func (e *exprAliased) Expr() q.Expr  { return e.expr }

func (e *exprAliased) ToQuery() q.Query {
	r := e.expr.ToQuery()
	return q.Query{
		fmt.Sprintf("%s AS %s", r.Query, e.alias),
		r.Args,
	}
}

func (e *exprAliased) ToSelectItem() q.SelectItem {
	item := e.expr.ToSelectItem()
	item.Alias = e.alias
	return item
}

type funcExpr struct {
	name string
	args []q.Queryable
}

func (f *funcExpr) FuncName() string { return f.name }

func (f *funcExpr) ToQuery() q.Query {
	qs := []string{}
	as := []interface{}{}
	for _, a := range f.args {
		qr := a.ToQuery()
		qs = append(qs, qr.Query)
		as = append(as, qr.Args...)
	}

	return q.Query{
		fmt.Sprintf("%s(%s)", f.name, strings.Join(qs, ", ")),
		as,
	}
}

func (f *funcExpr) ToSelectItem() q.SelectItem { return q.SelectItem{} }

type infixOp struct {
	left  q.Expr
	right q.Expr
	op    string
}

func (op *infixOp) ToQuery() q.Query {
	l := op.left.ToQuery()
	r := op.right.ToQuery()
	return q.Query{
		fmt.Sprintf("%s %s %s", l.Query, op.op, r.Query),
		append(l.Args, r.Args...),
	}
}

func (f *infixOp) ToSelectItem() q.SelectItem { return q.SelectItem{} }

type exprListExpr struct {
	qs []q.Queryable
}

func (el *exprListExpr) ToQuery() q.Query {
	ss := []string{}
	args := []interface{}{}
	for _, e := range el.qs {
		qr := e.ToQuery()
		ss = append(ss, qr.Query)
		args = append(args, qr.Args...)
	}
	return q.Query{
		strings.Join(ss, ", "),
		args,
	}
}

func (el *exprListExpr) ToSelectItem() q.SelectItem {
	panic("[INVALID] exprListExpr.ToSelectItem is called")
}

func (el *exprListExpr) Queryables() []q.Queryable {
	return el.qs
}
