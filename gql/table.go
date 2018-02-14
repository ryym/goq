package gql

type Table struct {
	name  string
	alias string
}

func NewTable(name, alias string) Table {
	return Table{name, alias}
}

func (t *Table) ApplyTable(q *Query, ctx DBContext) {
	q.query = append(q.query, ctx.QuoteIdent(t.name))
	if t.alias != "" {
		q.query = append(q.query, " AS ", ctx.QuoteIdent(t.alias))
	}
}

type DynmTable struct {
	Table
}

func (t *DynmTable) As(alias string) *DynmTable {
	return &DynmTable{NewTable(t.name, alias)}
}
