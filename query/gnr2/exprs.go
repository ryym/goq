package gnr2

import (
	"fmt"
	"strings"
)

type predExpr struct {
	Expr
}

func (p *predExpr) PredExpr() {}

type litExpr struct {
	val interface{}
	Ops
}

func (l *litExpr) init() *litExpr {
	l.Ops = Ops{l}
	return l
}

func (l *litExpr) Query() Query {
	return Query{"?", []interface{}{l.val}}
}

func (l *litExpr) SelectItem() SelectItem { return SelectItem{} }

type infixOp struct {
	left  Expr
	right Expr
	op    string
	Ops
}

func (op *infixOp) init() *infixOp {
	op.Ops = Ops{op}
	return op
}

func (op *infixOp) Query() Query {
	l := op.left.Query()
	r := op.right.Query()
	return Query{
		fmt.Sprintf("%s %s %s", l.Query, op.op, r.Query),
		append(l.Args, r.Args...),
	}
}

func (f *infixOp) SelectItem() SelectItem { return SelectItem{} }

type exprAliased struct {
	expr  Expr
	alias string
}

func (e *exprAliased) Alias() string { return e.alias }

func (e *exprAliased) Query() Query {
	r := e.expr.Query()
	return Query{
		fmt.Sprintf("%s AS %s", r.Query, e.alias),
		r.Args,
	}
}

type parensExpr struct {
	exp Expr
	Ops
}

func (p *parensExpr) init() *parensExpr {
	p.Ops = Ops{p}
	return p
}

func (p *parensExpr) Query() Query {
	qr := p.exp.Query()
	return Query{fmt.Sprintf("(%s)", qr.Query), qr.Args}
}

func (p *parensExpr) SelectItem() SelectItem { return p.exp.SelectItem() }

func (e *exprAliased) SelectItem() SelectItem {
	item := e.expr.SelectItem()
	item.Alias = e.alias
	return item
}

type logicalOp struct {
	op    string
	preds []PredExpr
	Ops
}

func (l *logicalOp) PredExpr() {}

func (l *logicalOp) init() *logicalOp {
	l.Ops = Ops{l}
	return l
}

func (l *logicalOp) Query() Query {
	if len(l.preds) == 0 {
		return Query{"", []interface{}{}}
	}

	pred := l.preds[0]
	for i := 1; i < len(l.preds); i++ {
		pred = &predExpr{(&infixOp{
			left:  pred,
			right: l.preds[i],
			op:    l.op,
		}).init()}
	}

	qr := pred.Query()
	return Query{
		fmt.Sprintf("(%s)", qr.Query),
		qr.Args,
	}
}

func (l *logicalOp) SelectItem() SelectItem { return SelectItem{} }

type exprListExpr struct {
	exps []Expr
}

func (el *exprListExpr) Query() Query {
	qs := []string{}
	args := []interface{}{}
	for _, e := range el.exps {
		qr := e.Query()
		qs = append(qs, qr.Query)
		args = append(args, qr.Args...)
	}
	return Query{
		strings.Join(qs, ", "),
		args,
	}
}

func (el *exprListExpr) SelectItem() SelectItem {
	panic("[INVALID] exprListExpr.SelectItem is called")
}

func (el *exprListExpr) Exprs() []Expr {
	return el.exps
}
