package gnr2

type Ops struct {
	expr Expr
}

func (op *Ops) Hello() Expr {
	return op.expr
}

func (op *Ops) As(alias string) ExprAliased {
	return &exprAliased{op.expr, alias}
}

func (op *Ops) Eq(v interface{}) PredExpr {
	return &predExpr{(&infixOp{
		left:  op.expr,
		right: lift(v),
		op:    "=",
	}).init()}
}

func (op *Ops) Add(v interface{}) Expr {
	return (&infixOp{
		left:  op.expr,
		right: lift(v),
		op:    "+",
	}).init()
}

// func (op *Ops) Mlt(v interface{}) Expr {
// 	return &Ops{&infixOp{op, lift(v), "*"}}
// }
