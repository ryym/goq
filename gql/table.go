package gql

type Table struct {
	name  string
	alias string
	cols  []*Column
}

func NewTable(name, alias string, cols []*Column) Table {
	return Table{name, alias, cols}
}

func (t *Table) ApplyTable(q *Query, ctx DBContext) {
	q.query = append(q.query, ctx.QuoteIdent(t.name))
	if t.alias != "" {
		q.query = append(q.query, " AS ", ctx.QuoteIdent(t.alias))
	}
}

func (t *Table) All() *ColumnListExpr {
	return NewColumnList(t.cols)
}

type DynmTable struct {
	table Table
}

func newDynmTable(name string) *DynmTable {
	return &DynmTable{NewTable(name, "", nil)}
}

func (t *DynmTable) ApplyTable(q *Query, ctx DBContext) {
	t.table.ApplyTable(q, ctx)
}

func (t *DynmTable) As(alias string) *DynmTable {
	return &DynmTable{NewTable(t.table.name, alias, nil)}
}
