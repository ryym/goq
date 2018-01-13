package gql

type queryExpr struct {
	exps   []Querier
	froms  []Table
	wheres []PredExpr
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

	// FROM
	if len(qe.froms) > 0 {
		q.query = append(q.query, " FROM ")
		lastIdx := len(qe.froms) - 1
		for i, t := range qe.froms {
			name := t.TableName()
			if alias := t.TableAlias(); alias != "" {
				name += " AS " + alias
			}
			q.query = append(q.query, name)
			if i < lastIdx {
				q.query = append(q.query, ", ")
			}
		}
	}

	// WHERE
	if len(qe.wheres) > 0 {
		q.query = append(q.query, " WHERE ")
		(&logicalOp{op: "AND", preds: qe.wheres}).Apply(q, ctx)
	}
}

func (qe *queryExpr) Selections() []Selection {
	items := make([]Selection, 0, len(qe.exps))
	for _, exp := range qe.exps {
		if cl, ok := exp.(ExprListExpr); ok {
			for _, e := range cl.Exprs() {
				items = append(items, e.Selection())
			}
		} else {
			items = append(items, exp.Selection())
		}
	}
	return items
}

func (qe *queryExpr) From(table Table, tables ...Table) Clauses {
	qe.froms = append(qe.froms, table)
	qe.froms = append(qe.froms, tables...)
	return qe
}

func (qe *queryExpr) Where(preds ...PredExpr) Clauses {
	qe.wheres = append(qe.wheres, preds...)
	return qe
}
