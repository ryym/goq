package goql

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

func (t *Table) All() *ColumnList {
	return NewColumnList(t.cols)
}

func (t *Table) Except(excepts ...*Column) *ColumnList {
	if len(excepts) == 0 {
		return t.All()
	}

	var cols []*Column
	for _, col := range t.cols {
		isTarget := true
		for _, ecol := range excepts {
			if col == ecol {
				isTarget = false
				break
			}
		}
		if isTarget {
			cols = append(cols, col)
		}
	}
	return NewColumnList(cols)
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
