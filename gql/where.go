package gql

type Where struct {
	preds []PredExpr
}

func (w *Where) add(preds []PredExpr) {
	w.preds = append(w.preds, preds...)
}

func (w *Where) Apply(q *Query, ctx DBContext) {
	if len(w.preds) > 0 {
		q.query = append(q.query, " WHERE ")
		(&logicalOp{op: "AND", preds: w.preds}).Apply(q, ctx)
	}
}
