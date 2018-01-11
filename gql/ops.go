package gql

func lift(v interface{}) Expr {
	exp, ok := v.(Expr)
	if ok {
		return exp
	}
	return (&litExpr{val: v}).init()
}

type ops struct {
	expr Expr
}

func (o *ops) As(alias string) Aliased {
	return &aliased{o.expr, alias}
}

func (o *ops) Eq(v interface{}) PredExpr {
	return &predExpr{(&infixOp{
		left:  o.expr,
		right: lift(v),
		op:    "=",
	}).init()}
}

func (o *ops) Neq(v interface{}) PredExpr {
	return &predExpr{(&infixOp{
		left:  o.expr,
		right: lift(v),
		op:    "<>",
	}).init()}
}

func (o *ops) Gt(v interface{}) PredExpr {
	return &predExpr{(&infixOp{
		left:  o.expr,
		right: lift(v),
		op:    ">",
	}).init()}
}

func (o *ops) Gte(v interface{}) PredExpr {
	return &predExpr{(&infixOp{
		left:  o.expr,
		right: lift(v),
		op:    ">=",
	}).init()}
}

func (o *ops) Lt(v interface{}) PredExpr {
	return &predExpr{(&infixOp{
		left:  o.expr,
		right: lift(v),
		op:    "<",
	}).init()}
}

func (o *ops) Lte(v interface{}) PredExpr {
	return &predExpr{(&infixOp{
		left:  o.expr,
		right: lift(v),
		op:    "<=",
	}).init()}
}

func (o *ops) Like(s string) PredExpr {
	return &predExpr{(&infixOp{
		left:  o.expr,
		right: lift(s),
		op:    "LIKE",
	}).init()}
}

func (p *ops) Between(start interface{}, end interface{}) PredExpr {
	return &predExpr{(&betweenOp{
		start: lift(start),
		end:   lift(end),
	})}
}

func (o *ops) IsNull() PredExpr {
	return &predExpr{(&sufixOp{
		val: o.expr,
		op:  "IS NULL",
	})}
}

func (o *ops) IsNotNull() PredExpr {
	return &predExpr{(&sufixOp{
		val: o.expr,
		op:  "IS NOT NULL",
	})}
}
