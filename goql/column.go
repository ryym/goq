package goql

type ColumnMaker struct {
	structName string
	tableName  string
	tableAlias string
}

func NewColumnMaker(structName, tableName string) *ColumnMaker {
	return &ColumnMaker{structName: structName, tableName: tableName}
}

func (m *ColumnMaker) As(alias string) *ColumnMaker {
	m.tableAlias = alias
	return m
}

func (m *ColumnMaker) Col(fieldName, name string) *ColumnBuilder {
	col := (&Column{
		tableName:  m.tableName,
		tableAlias: m.tableAlias,
		structName: m.structName,
		name:       name,
		fieldName:  fieldName,
		meta:       ColumnMeta{},
	}).init()
	return &ColumnBuilder{col}
}

type ColumnBuilder struct {
	col *Column
}

func (cb *ColumnBuilder) PK() *ColumnBuilder {
	cb.col.meta.PK = true
	return cb
}

func (cb *ColumnBuilder) Bld() *Column {
	return cb.col
}

type ColumnMeta struct {
	PK bool
}

type Column struct {
	tableName  string
	tableAlias string
	structName string
	name       string
	fieldName  string
	meta       ColumnMeta
	ops
}

func (c *Column) init() *Column {
	c.ops = ops{c}
	return c
}

func (c *Column) TableName() string  { return c.tableName }
func (c *Column) TableAlias() string { return c.tableAlias }
func (c *Column) StructName() string { return c.structName }
func (c *Column) ColumnName() string { return c.name }
func (c *Column) FieldName() string  { return c.fieldName }
func (c *Column) Meta() *ColumnMeta  { return &c.meta }

func (c *Column) Apply(q *Query, ctx DBContext) {
	table := c.tableAlias
	if table == "" {
		table = c.tableName
	}
	if table != "" {
		q.query = append(q.query, ctx.QuoteIdent(table), ".")
	}
	q.query = append(q.query, ctx.QuoteIdent(c.name))
}

func (c *Column) Selection() Selection {
	return Selection{
		ColumnName: c.name,
		TableName:  c.tableName,
		TableAlias: c.tableAlias,
		StructName: c.structName,
		FieldName:  c.fieldName,
	}
}
