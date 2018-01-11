package gnr2

import "reflect"

type Goq struct{}

func (g *Goq) Parens(exp Querier) Expr {
	return (&parensExpr{exp: exp}).init()
}

func (g *Goq) And(preds ...PredExpr) PredExpr {
	return &predExpr{(&logicalOp{
		op:    "AND",
		preds: preds,
	}).init()}
}

func (g *Goq) Select(exps ...Querier) SelectClause {
	return &selectClause{exps}
}

func (g *Goq) InnerJoin(table Table) *Joinner {
	return &Joinner{table}
}

type Joinner struct {
	table Table
}

func (j *Joinner) On(exp PredExpr) JoinOn {
	return JoinOn{j.table, exp, JOIN_INNER}
}

func (g *Goq) ToSlice(slice interface{}) *GSliceCollector {
	sliceV := reflect.ValueOf(slice).Elem()
	return &GSliceCollector{
		itemType: sliceV.Type().Elem(),
		slice:    sliceV,
	}
}
