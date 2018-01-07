package gnr

import q "github.com/ryym/goq/query"

type predExpr struct {
	*Ops
}

func (p *predExpr) PredExpr() {}

// 各種 expr の struct は Ops で包む事で q.Expr を満たし、
// 演算可能になる。

type Ops struct {
	q.Queryable
}

func (op *Ops) As(alias string) q.ExprAliased {
	return &exprAliased{op, alias}
}

func (op *Ops) Eq(v interface{}) q.PredExpr {
	return &predExpr{&Ops{&infixOp{op, lift(v), "="}}}
}

func (op *Ops) Add(v interface{}) q.Expr {
	return &Ops{&infixOp{op, lift(v), "+"}}
}

func (op *Ops) Mlt(v interface{}) q.Expr {
	return &Ops{&infixOp{op, lift(v), "*"}}
}
