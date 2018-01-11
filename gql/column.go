package gql

import "fmt"

func NewColumnMaker(tableName, structName string) ColumnMaker {
	return ColumnMaker{tableName, structName}
}

type ColumnMaker struct {
	tableName  string
	structName string
}

func (m *ColumnMaker) Col(fieldName, name string) *column {
	col := column{
		tableName:  m.tableName,
		tableAlias: "",
		structName: m.structName,
		name:       name,
		fieldName:  fieldName,
	}.init()
	return &col
}

type column struct {
	tableName  string
	tableAlias string
	structName string
	name       string
	fieldName  string
	ops
}

func (c column) init() column {
	c.ops = ops{&c}
	return c
}

func (c *column) TableName() string  { return c.tableName }
func (c *column) TableAlias() string { return c.tableAlias }
func (c *column) StructName() string { return c.structName }
func (c *column) ColumnName() string { return c.name }
func (c *column) FieldName() string  { return c.fieldName }

func (c *column) Query() Query {
	table := c.tableAlias
	if table == "" {
		table = c.tableName
	}
	return Query{
		fmt.Sprintf("%s.%s", table, c.name),
		[]interface{}{},
	}
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
