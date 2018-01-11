package gql

import "fmt"

type aliased struct {
	expr  Expr
	alias string
}

func (a *aliased) Alias() string { return a.alias }

func (a *aliased) Query() Query {
	qr := a.expr.Query()
	return Query{
		fmt.Sprintf("%s AS %s", qr.Query, a.alias),
		qr.Args,
	}
}

func (a *aliased) Selection() Selection { return Selection{} }

type litExpr struct {
	val interface{}
	ops
}

func (l *litExpr) init() *litExpr {
	l.ops = ops{l}
	return l
}

// TODO: Add no placeholder version
func (l *litExpr) Query() Query {
	return Query{"?", []interface{}{l.val}}
}

func (l *litExpr) Selection() Selection { return Selection{} }

type predExpr struct {
	Expr
}

func (p *predExpr) ImplPredExpr() {}

type rawExpr struct {
	sql string
	ops
}

func (r *rawExpr) init() *rawExpr {
	r.ops = ops{r}
	return r
}

func (r *rawExpr) Query() Query {
	return Query{r.sql, []interface{}{}}
}

func (r *rawExpr) Selection() Selection { return Selection{} }

type parensExpr struct {
	exp Expr
	ops
}

func (p *parensExpr) init() *parensExpr {
	p.ops = ops{p}
	return p
}

func (p *parensExpr) Query() Query {
	qr := p.exp.Query()
	return Query{
		fmt.Sprintf("(%s)", qr.Query),
		qr.Args,
	}
}

func (p *parensExpr) Selection() Selection { return p.exp.Selection() }
