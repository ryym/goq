package gql

type Builder struct{}

func (b *Builder) Query(exp Querier) Query {
	ctx := &postgresCtx{} // TODO: Change dynamically.
	q := Query{}
	exp.Apply(&q, ctx)
	return q
}

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
		op: "NOT", val: pred,
	}).init()}
}

func (b *Builder) Func(name string, args ...interface{}) Expr {
	expArgs := make([]Expr, len(args))
	for i, a := range args {
		expArgs[i] = lift(a)
	}
	return (&funcExpr{name: name, args: expArgs}).init()
}

func (b *Builder) Count(exp Expr) Expr {
	return b.Func("COUNT", exp)
}

func (b *Builder) Sum(exp Expr) Expr {
	return b.Func("SUM", exp)
}

func (b *Builder) Min(exp Expr) Expr {
	return b.Func("MIN", exp)
}

func (b *Builder) Max(exp Expr) Expr {
	return b.Func("MAX", exp)
}

func (b *Builder) Avg(exp Expr) Expr {
	return b.Func("AVG", exp)
}

func (b *Builder) Coalesce(exp Expr, alt interface{}) Expr {
	return b.Func("COALESCE", exp, lift(alt))
}

func (b *Builder) Concat(exps ...interface{}) Expr {
	return b.Func("CONCAT", exps...)
}

func (b *Builder) Exists(query QueryExpr) PredExpr {
	return &predExpr{(&existsExpr{query: query}).init()}
}

func (b *Builder) Select(exps ...Querier) SelectClause {
	return (&queryExpr{exps: exps}).init()
}

// TODO: Accept a sub query.
func (b *Builder) InnerJoin(table Table) *JoinClause {
	return &JoinClause{JOIN_INNER, table}
}

func (b *Builder) LeftJoin(table Table) *JoinClause {
	return &JoinClause{JOIN_LEFT, table}
}

func (b *Builder) RightJoin(table Table) *JoinClause {
	return &JoinClause{JOIN_RIGHT, table}
}

func (b *Builder) FullJoin(table Table) *JoinClause {
	return &JoinClause{JOIN_FULL, table}
}

func (b *Builder) Case(cases ...*WhenExpr) *CaseExpr {
	return (&CaseExpr{cases: cases}).init()
}

type CaseOfExpr func(cases ...*WhenExpr) *CaseExpr

func (b *Builder) CaseOf(val Expr) CaseOfExpr {
	return func(cases ...*WhenExpr) *CaseExpr {
		return (&CaseExpr{val: val, cases: cases}).init()
	}
}

func (b *Builder) When(when interface{}) *WhenExpr {
	return &WhenExpr{when: lift(when)}
}
