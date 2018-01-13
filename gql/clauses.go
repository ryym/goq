package gql

type queryExpr struct {
	exps []Querier
	ops
}

func (qe *queryExpr) init() *queryExpr {
	qe.ops = ops{qe}
	return qe
}

func (q *queryExpr) Selection() Selection { return Selection{} }

func (qe *queryExpr) Apply(q *Query, ctx DBContext) {
	q.query = append(q.query, "SELECT ")
	if len(qe.exps) == 0 {
		return // XXX: Should return an error?
	}
	qe.exps[0].Apply(q, ctx)
	for i := 1; i < len(qe.exps); i++ {
		q.query = append(q.query, ", ")
		qe.exps[i].Apply(q, ctx)
	}
}

func (qe *queryExpr) Selections() []Selection {
	items := make([]Selection, 0, len(qe.exps))
	for _, exp := range qe.exps {
		items = append(items, exp.Selection())
	}
	return items
}
