package gql

func NewColumnMaker(tableName, structName string) ColumnMaker {
	return ColumnMaker{tableName, structName}
}

type ColumnMaker struct {
	tableName  string
	structName string
}

func (m *ColumnMaker) Col(fieldName, name string) *column {
	return (&column{
		tableName:  m.tableName,
		tableAlias: "",
		structName: m.structName,
		name:       name,
		fieldName:  fieldName,
	}).init()
}

type column struct {
	tableName  string
	tableAlias string
	structName string
	name       string
	fieldName  string
	ops
}

func (c *column) init() *column {
	c.ops = ops{c}
	return c
}

func (c *column) TableName() string  { return c.tableName }
func (c *column) TableAlias() string { return c.tableAlias }
func (c *column) StructName() string { return c.structName }
func (c *column) ColumnName() string { return c.name }
func (c *column) FieldName() string  { return c.fieldName }

func (c *column) Apply(q *Query, ctx DBContext) {
	table := c.tableAlias
	if table == "" {
		table = c.tableName
	}
	q.query = append(
		q.query,
		ctx.QuoteIdent(table),
		".",
		ctx.QuoteIdent(c.name),
	)
}

func (c *column) Selection() Selection {
	return Selection{
		ColumnName: c.name,
		TableName:  c.tableName,
		TableAlias: c.tableAlias,
		StructName: c.structName,
		FieldName:  c.fieldName,
	}
}
