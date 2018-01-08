package gnr

import (
	"fmt"
	"strings"

	q "github.com/ryym/goq/query"
)

type selectClause struct {
	exps []q.Queryable
}

func (sc *selectClause) From(table q.Table, tables ...q.Table) q.Clauses {
	froms := make([]q.Table, len(tables)+1)
	froms[0] = table
	for i, t := range tables {
		froms[i+1] = t
	}
	return &clauses{selectCls: sc, froms: froms}
}

func (sc *selectClause) ToSelectItem() q.SelectItem { return q.SelectItem{} }

func (sc *selectClause) ToQuery() q.Query {
	selects := make([]string, len(sc.exps))
	args := []interface{}{}
	for i, exp := range sc.exps {
		qr := exp.ToQuery()
		selects[i] = qr.Query
		args = append(args, qr.Args...)
	}
	return q.Query{
		"SELECT " + strings.Join(selects, ", "),
		args,
	}
}

func (sc *selectClause) GetSelects() []q.SelectItem {
	items := make([]q.SelectItem, 0, len(sc.exps))
	for _, exp := range sc.exps {
		if cl, ok := exp.(q.ExprListExpr); ok {
			for _, qb := range cl.Queryables() {
				items = append(items, qb.ToSelectItem())
			}
		} else {
			items = append(items, exp.ToSelectItem())
		}
	}
	return items
}

type clauses struct {
	selectCls *selectClause
	froms     []q.Table
	wheres    []q.PredExpr
	joins     []q.JoinOn
	orders    []q.Expr
	limit     int
	offset    int
}

func (cl *clauses) GetSelects() []q.SelectItem {
	return cl.selectCls.GetSelects()
}

func (cl *clauses) ToQuery() q.Query {
	qr := cl.selectCls.ToQuery()

	qs := []string{}
	qs = append(qs, qr.Query)

	args := qr.Args[:]

	// FROM
	if len(cl.froms) > 0 {
		ts := make([]string, len(cl.froms))
		for i, t := range cl.froms {
			ts[i] = t.TableName()
			if alias := t.Alias(); alias != "" {
				ts[i] += fmt.Sprintf(" AS %s", alias)
			}
		}
		qs = append(qs, "FROM "+strings.Join(ts, ", "))
	}

	// JOIN
	for _, j := range cl.joins {
		table := j.Table.TableName()
		if alias := j.Table.Alias(); alias != "" {
			table += fmt.Sprintf(" AS %s", alias)
		}
		jqr := j.On.ToQuery()
		args = append(args, jqr.Args...)
		qs = append(qs, fmt.Sprintf(
			"%s JOIN %s ON %s",
			j.Kind,
			table,
			jqr.Query,
		))
	}

	// WHERE
	if len(cl.wheres) > 0 {
		qr = (&logicalOp{"AND", cl.wheres}).ToQuery()
		qs = append(qs, "WHERE "+qr.Query)
		args = append(args, qr.Args...)
	}

	return q.Query{strings.Join(qs, " "), args}
}

func (cl *clauses) ToSelectItem() q.SelectItem {
	return q.SelectItem{}
}

func (cl *clauses) Where(exps ...q.PredExpr) q.Clauses {
	cl.wheres = append(cl.wheres, exps...)
	return cl
}

func (cl *clauses) Joins(joins ...q.JoinOn) q.Clauses {
	cl.joins = append(cl.joins, joins...)
	return cl
}

func (cl *clauses) GroupBy(exps ...q.Queryable) q.GroupQuery {
	return nil
}

func (cl *clauses) Having(exps ...q.PredExpr) q.GroupQuery {
	return nil
}

func (cl *clauses) OrderBy(exps ...q.Queryable) q.ExtraClauses {
	return nil
}

func (cl *clauses) Limit(n int) q.ExtraClauses {
	return nil
}

func (cl *clauses) Offset(n int) q.ExtraClauses {
	return nil
}
