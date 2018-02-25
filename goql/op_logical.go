package goql

type logicalOp struct {
	op    string
	preds []PredExpr
	ops
}

func (l *logicalOp) ImplPredExpr() {}

func (l *logicalOp) init() *logicalOp {
	l.ops = ops{l}
	return l
}

func (l *logicalOp) Apply(q *Query, ctx DBContext) {
	if len(l.preds) == 0 {
		return
	}
	pred := concatPreds(l.preds, l.op)
	q.query = append(q.query, "(")
	pred.Apply(q, ctx)
	q.query = append(q.query, ")")
}

func (l *logicalOp) Selection() Selection { return Selection{} }

func concatPreds(preds []PredExpr, op string) PredExpr {
	pred := preds[0]
	for i := 1; i < len(preds); i++ {
		pred = &predExpr{(&infixOp{
			left:  pred,
			right: preds[i],
			op:    op,
		}).init()}
	}
	return pred
}
