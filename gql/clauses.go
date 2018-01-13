package gql

import "fmt"

type queryExpr struct {
	exps   []Querier
	froms  []Table
	joins  []JoinOn
	wheres []PredExpr
	groups []Expr
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
			table := t.TableName()
			if alias := t.TableAlias(); alias != "" {
				table += " AS " + alias
			}
			q.query = append(q.query, table)
			if i < lastIdx {
				q.query = append(q.query, ", ")
			}
		}
	}

	// JOIN
	for _, j := range qe.joins {
		table := j.Table.TableName()
		if alias := j.Table.TableAlias(); alias != "" {
			table += " AS " + alias
		}
		q.query = append(q.query, fmt.Sprintf(" %s JOIN %s ON ", j.Type, table))
		j.On.Apply(q, ctx)
	}

	// WHERE
	if len(qe.wheres) > 0 {
		q.query = append(q.query, " WHERE ")
		(&logicalOp{op: "AND", preds: qe.wheres}).Apply(q, ctx)
	}

	// GROUP BY
	if len(qe.groups) > 0 {
		q.query = append(q.query, " GROUP BY ")
		qe.groups[0].Apply(q, ctx)
		for i := 1; i < len(qe.groups); i++ {
			q.query = append(q.query, ", ")
			qe.groups[i].Apply(q, ctx)
		}
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

func (qe *queryExpr) Joins(joins ...JoinOn) Clauses {
	qe.joins = append(qe.joins, joins...)
	return qe
}

func (qe *queryExpr) Where(preds ...PredExpr) Clauses {
	qe.wheres = append(qe.wheres, preds...)
	return qe
}

func (qe *queryExpr) GroupBy(exps ...Expr) GroupByClause {
	qe.groups = append(qe.groups, exps...)
	return qe
}
