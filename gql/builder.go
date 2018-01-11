package gql

type Builder struct{}

func (b *Builder) Var(v interface{}) Expr {
	return (&litExpr{val: v}).init()
}

func (b *Builder) Raw(sql string) Expr {
	return (&rawExpr{sql: sql}).init()
}

func (b *Builder) Parens(exp Expr) Expr {
	return (&parensExpr{exp: exp}).init()
}

func (b *Builder) And(preds ...PredExpr) PredExpr {
	return (&logicalOp{op: "AND", preds: preds}).init()
}

func (b *Builder) Or(preds ...PredExpr) PredExpr {
	return (&logicalOp{op: "OR", preds: preds}).init()
}

func (b *Builder) Not(pred PredExpr) PredExpr {
	return &predExpr{(&prefixOp{
		op: "NOT ", val: pred,
	}).init()}
}

func (b *Builder) Func(name string, args ...interface{}) Expr {
	expArgs := make([]Expr, len(args))
	for i, a := range args {
		expArgs[i] = lift(a)
	}
	return (&funcExpr{name: name, args: expArgs}).init()
}
