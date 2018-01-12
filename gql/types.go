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

	Eq(v interface{}) PredExpr
	Neq(v interface{}) PredExpr
	Gt(v interface{}) PredExpr
	Gte(v interface{}) PredExpr
	Lt(v interface{}) PredExpr
	Lte(v interface{}) PredExpr
	Like(s string) PredExpr
	Between(start interface{}, end interface{}) PredExpr
	IsNull() PredExpr
	IsNotNull() PredExpr

	Add(v interface{}) Expr
	Sbt(v interface{}) Expr
	Mlt(v interface{}) Expr
	Dvd(v interface{}) Expr
	Concat(s interface{}) Expr
}

// PredExpr represents this expression is a predicate.
type PredExpr interface {
	Expr
	ImplPredExpr()
}

type Column interface {
	Expr
	TableName() string
	TableAlias() string
	StructName() string
	ColumnName() string
	FieldName() string
}
