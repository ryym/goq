package gql

func AllCols(cols []Column) ExprListExpr {
	exps := make([]Expr, len(cols))
	for i, c := range cols {
		exps[i] = c
	}
	return &exprListExpr{exps}
}

type Table struct {
	name  string
	alias string
}

func NewTable(name, alias string) Table {
	return Table{name, alias}
}

func (t *Table) TableName() string  { return t.name }
func (t *Table) TableAlias() string { return t.alias }

func (t *Table) ApplyTable(q *Query, ctx DBContext) {
	q.query = append(q.query, ctx.QuoteIdent(t.name))
	if t.alias != "" {
		q.query = append(q.query, " AS ", ctx.QuoteIdent(t.alias))
	}
}

type DynmTable struct {
	name  string
	alias string
}

func (t *DynmTable) As(alias string) *DynmTable {
	t.alias = alias
	return t
}

func (t *DynmTable) ApplyTable(q *Query, ctx DBContext) {
	q.query = append(q.query, ctx.QuoteIdent(t.name))
	if t.alias != "" {
		q.query = append(q.query, " AS ", ctx.QuoteIdent(t.alias))
	}
}
