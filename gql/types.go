package gql

type Query struct {
	Query string
	Args  []interface{}
}

type Selection struct {
	Alias      string
	ColumnName string
	TableName  string
	TableAlias string
	StructName string
	FieldName  string
}

type Querier interface {
	Query() Query
	Selection() Selection
}

type Aliased interface {
	Querier
	Alias() string
}

type Expr interface {
	Querier
	As(alias string) Aliased
	// TODO: More operators
}
