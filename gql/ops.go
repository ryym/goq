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

func (o *ops) Between(start interface{}, end interface{}) PredExpr {
	return &predExpr{(&betweenOp{
		val:   o.expr,
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

func (o *ops) Add(v interface{}) Expr {
	return (&infixOp{
		left:  o.expr,
		right: lift(v),
		op:    "+",
	}).init()
}

func (o *ops) Sbt(v interface{}) Expr {
	return (&infixOp{
		left:  o.expr,
		right: lift(v),
		op:    "-",
	}).init()
}

func (o *ops) Mlt(v interface{}) Expr {
	return (&infixOp{
		left:  o.expr,
		right: lift(v),
		op:    "*",
	}).init()
}

func (o *ops) Dvd(v interface{}) Expr {
	return (&infixOp{
		left:  o.expr,
		right: lift(v),
		op:    "/",
	}).init()
}

func (o *ops) Concat(v interface{}) Expr {
	return (&infixOp{
		left:  o.expr,
		right: lift(v),
		op:    "||",
	}).init()
}

func (o *ops) In(vals ...interface{}) PredExpr {
	exps := make([]Expr, len(vals))
	for i, v := range vals {
		exps[i] = lift(v)
	}
	return &predExpr{(&inExpr{
		val:  o.expr,
		exps: exps,
	}).init()}
}
