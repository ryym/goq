package gql

import "github.com/ryym/goq/dialect"

func NewBuilder(dl dialect.Dialect) *Builder {
	return &Builder{dl}
}

type Builder struct {
	ctx DBContext
}

func (b *Builder) Query(exp QueryApplier) Query {
	q := Query{}
	exp.Apply(&q, b.ctx)
	return q
}

func (b *Builder) Var(v interface{}) AnonExpr {
	return (&litExpr{val: v}).init()
}

func (b *Builder) VarT(v interface{}, typ string) AnonExpr {
	return (&litExpr{val: v, typ: typ}).init()
}

// XXX: This should be able to use everywhere.
func (b *Builder) Raw(sql string) AnonExpr {
	return (&rawExpr{sql: sql}).init()
}

func (b *Builder) Parens(exp Expr) AnonExpr {
	return (&parensExpr{exp: exp}).init()
}

// XXX: This should be used only in ORDER BY?
func (b *Builder) Name(name string) *nameExpr {
	return (&nameExpr{name: name}).init()
}

func (b *Builder) Col(table, col string) *Column {
	return (&Column{
		tableName:  "",
		tableAlias: table,
		structName: "",
		name:       col,
		fieldName:  "",
	}).init()
}

func (b *Builder) Table(name string) *DynmTable {
	return &DynmTable{NewTable(name, "")}
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

func (b *Builder) Func(name string, args ...interface{}) AnonExpr {
	expArgs := make([]Expr, len(args))
	for i, a := range args {
		expArgs[i] = lift(a)
	}
	return (&funcExpr{name: name, args: expArgs}).init()
}

func (b *Builder) Count(exp Expr) AnonExpr {
	return b.Func("COUNT", exp)
}

func (b *Builder) Sum(exp Expr) AnonExpr {
	return b.Func("SUM", exp)
}

func (b *Builder) Min(exp Expr) AnonExpr {
	return b.Func("MIN", exp)
}

func (b *Builder) Max(exp Expr) AnonExpr {
	return b.Func("MAX", exp)
}

func (b *Builder) Avg(exp Expr) AnonExpr {
	return b.Func("AVG", exp)
}

func (b *Builder) Coalesce(exp Expr, alt interface{}) AnonExpr {
	return b.Func("COALESCE", exp, lift(alt))
}

func (b *Builder) Concat(exps ...interface{}) AnonExpr {
	return b.Func("CONCAT", exps...)
}

func (b *Builder) Exists(query QueryExpr) PredExpr {
	return &predExpr{(&existsExpr{query: query}).init()}
}

func (b *Builder) Select(exps ...Querier) SelectClause {
	return (&queryExpr{exps: exps, ctx: b.ctx}).init()
}

// TODO: Accept a sub query.
func (b *Builder) InnerJoin(table TableLike) *JoinClause {
	return &JoinClause{JOIN_INNER, table}
}

func (b *Builder) LeftJoin(table TableLike) *JoinClause {
	return &JoinClause{JOIN_LEFT, table}
}

func (b *Builder) RightJoin(table TableLike) *JoinClause {
	return &JoinClause{JOIN_RIGHT, table}
}

func (b *Builder) FullJoin(table TableLike) *JoinClause {
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
