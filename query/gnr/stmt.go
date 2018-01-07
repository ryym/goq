package gnr

import (
	"strings"

	q "github.com/ryym/goq/query"
)

type selectClause struct {
	exps []q.Queryable
}

func (sc *selectClause) From(table q.TableBase, tables ...q.TableBase) q.Clauses {
	return nil
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
