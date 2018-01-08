package gnr2

type Goq struct{}

func (g *Goq) Parens(exp Expr) Expr {
	return (&parensExpr{exp: exp}).init()
}

func (g *Goq) And(preds ...PredExpr) PredExpr {
	return &predExpr{(&logicalOp{
		op:    "AND",
		preds: preds,
	}).init()}
}

func (g *Goq) Select(exps ...Expr) SelectClause {
	return &selectClause{exps}
}
