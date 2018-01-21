package gql

import (
	"fmt"
	"strconv"
)

type queryExpr struct {
	exps    []Querier
	froms   []TableLike
	joins   []*JoinDef
	wheres  []PredExpr
	groups  []Expr
	havings []PredExpr
	orders  []Expr
	limit   int
	offset  int
	ctx     DBContext
	ops
}

func (qe *queryExpr) init() *queryExpr {
	qe.ops = ops{qe}
	return qe
}

func (qe *queryExpr) Selection() Selection { return Selection{} }

func (qe *queryExpr) Construct() Query {
	q := Query{}
	qe.Apply(&q, qe.ctx)
	return q
}

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
			t.ApplyTable(q, ctx)
			if i < lastIdx {
				q.query = append(q.query, ", ")
			}
		}
	}

	// JOIN
	for _, j := range qe.joins {
		q.query = append(q.query, fmt.Sprintf(" %s JOIN ", j.Type))
		j.Table.ApplyTable(q, ctx)
		q.query = append(q.query, " ON ")
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

	// HAVING
	if len(qe.havings) > 0 {
		q.query = append(q.query, " HAVING ")
		(&logicalOp{op: "AND", preds: qe.havings}).Apply(q, ctx)
	}

	// ORDER BY
	if len(qe.orders) > 0 {
		q.query = append(q.query, " ORDER BY ")
		qe.orders[0].Apply(q, ctx)
		for i := 1; i < len(qe.orders); i++ {
			q.query = append(q.query, ", ")
			qe.orders[i].Apply(q, ctx)
		}
	}

	// LIMIT
	if qe.limit > 0 {
		q.query = append(q.query, " LIMIT ", strconv.Itoa(qe.limit))
	}

	// OFFSET
	if qe.offset > 0 {
		q.query = append(q.query, " OFFSET ", strconv.Itoa(qe.offset))
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

func (qe *queryExpr) From(table TableLike, tables ...TableLike) Clauses {
	qe.froms = append(qe.froms, table)
	qe.froms = append(qe.froms, tables...)
	return qe
}

func (qe *queryExpr) Joins(definers ...JoinDefiner) Clauses {
	for _, def := range definers {
		qe.joins = append(qe.joins, def.joinDef())
	}
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

func (qe *queryExpr) Having(preds ...PredExpr) GroupByClause {
	qe.havings = append(qe.havings, preds...)
	return qe
}

func (qe *queryExpr) OrderBy(exps ...Expr) QueryExpr {
	qe.orders = append(qe.orders, exps...)
	return qe
}

func (qe *queryExpr) Limit(n int) QueryExpr {
	qe.limit = n
	return qe
}

func (qe *queryExpr) Offset(n int) QueryExpr {
	qe.offset = n
	return qe
}

func (qe *queryExpr) As(alias string) QueryTable {
	return &aliasedQuery{aliased{
		(&parensExpr{exp: qe}).init(),
		alias,
	}}
}

type aliasedQuery struct {
	aliased
}

func (aq *aliasedQuery) ApplyTable(q *Query, ctx DBContext) {
	aq.aliased.Apply(q, ctx)
}
