package gnr

import q "github.com/ryym/goq/query"

type GoqGnr struct{}

func (g *GoqGnr) Select(exps ...q.Queryable) q.SelectClause {
	return &selectClause{exps}
}

func (g *GoqGnr) Val(v interface{}) q.Expr {
	return lift(v)
}

func (g *GoqGnr) Col(table string, name string) q.Expr {
	return &Ops{&columnExpr{name: name, tableName: table}}
}

func (g *GoqGnr) And(exps ...q.PredExpr) q.PredExpr {
	return &predExpr{&Ops{&logicalOp{"AND", exps}}}
}

func (g *GoqGnr) Parens(exp q.Expr) q.Expr {
	return &Ops{&parensExpr{exp}}
}

func (g *GoqGnr) Func(name string, args ...interface{}) q.Expr {
	expArgs := make([]q.Queryable, len(args))
	for i, a := range args {
		expArgs[i] = lift(a)
	}
	return &Ops{&funcExpr{name, expArgs}}
}

func (g *GoqGnr) InnerJoin(table q.Table) *NormalJoinner {
	return &NormalJoinner{table}
}

type NormalJoinner struct {
	table q.Table
}

func (j *NormalJoinner) On(exp q.PredExpr) q.JoinOn {
	return q.JoinOn{j.table, exp, q.JOIN_INNER}
}
