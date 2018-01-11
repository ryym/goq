package prot

import (
	"fmt"
	"strings"
)

type selectClause struct {
	exps []Querier
}

func (sc *selectClause) From(table Table, tables ...Table) Clauses {
	froms := make([]Table, len(tables)+1)
	froms[0] = table
	for i, t := range tables {
		froms[i+1] = t
	}
	return &clauses{selectCls: sc, froms: froms}
}

func (sc *selectClause) SelectItem() SelectItem { return SelectItem{} }

func (sc *selectClause) Query() Query {
	qs := make([]string, len(sc.exps))
	args := []interface{}{}
	for i, exp := range sc.exps {
		qr := exp.Query()
		qs[i] = qr.Query
		args = append(args, qr.Args...)
	}
	return Query{
		"SELECT " + strings.Join(qs, ", "),
		args,
	}
}

func (sc *selectClause) GetSelects() []SelectItem {
	items := make([]SelectItem, 0, len(sc.exps))
	for _, exp := range sc.exps {
		if cl, ok := exp.(ExprListExpr); ok {
			for _, e := range cl.Exprs() {
				items = append(items, e.SelectItem())
			}
		} else {
			items = append(items, exp.SelectItem())
		}
	}
	return items
}

type clauses struct {
	selectCls *selectClause
	froms     []Table
	wheres    []PredExpr
	joins     []JoinOn
	orders    []Expr
	limit     int
	offset    int
}

func (cl *clauses) GetSelects() []SelectItem {
	return cl.selectCls.GetSelects()
}

func (cl *clauses) Query() Query {
	qr := cl.selectCls.Query()

	qs := []string{}
	qs = append(qs, qr.Query)

	args := qr.Args[:]

	// FROM
	if len(cl.froms) > 0 {
		ts := make([]string, len(cl.froms))
		for i, t := range cl.froms {
			ts[i] = t.TableName()
			if alias := t.TableAlias(); alias != "" {
				ts[i] += fmt.Sprintf(" AS %s", alias)
			}
		}
		qs = append(qs, "FROM "+strings.Join(ts, ", "))
	}

	// JOIN
	for _, j := range cl.joins {
		table := j.Table.TableName()
		if alias := j.Table.TableAlias(); alias != "" {
			table += fmt.Sprintf(" AS %s", alias)
		}
		jqr := j.On.Query()
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
		qr = (&logicalOp{op: "AND", preds: cl.wheres}).Query()
		qs = append(qs, "WHERE "+qr.Query)
		args = append(args, qr.Args...)
	}

	// LIMIT
	if cl.limit > 0 {
		qs = append(qs, fmt.Sprintf("LIMIT %d", cl.limit))
	}

	return Query{strings.Join(qs, " "), args}
}

func (cl *clauses) SelectItem() SelectItem { return SelectItem{} }

func (cl *clauses) Where(exps ...PredExpr) Clauses {
	cl.wheres = append(cl.wheres, exps...)
	return cl
}

func (cl *clauses) Joins(joins ...JoinOn) Clauses {
	cl.joins = append(cl.joins, joins...)
	return cl
}

// TODO

func (cl *clauses) GroupBy(exps ...Expr) GroupQuery {
	return nil
}

func (cl *clauses) Having(exps ...PredExpr) GroupQuery {
	return nil
}

func (cl *clauses) OrderBy(exps ...Expr) ExtraClauses {
	return nil
}

func (cl *clauses) Limit(n int) ExtraClauses {
	cl.limit = n
	return cl
}

func (cl *clauses) Offset(n int) ExtraClauses {
	return nil
}
