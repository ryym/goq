package gql

import (
	"fmt"
	"strings"
)

type Query struct {
	query []string
	args  []interface{}
}

func (q *Query) Query() string {
	return strings.Join(q.query, "")
}

func (q *Query) Args() []interface{} {
	return q.args
}

func (q Query) String() string {
	return fmt.Sprintf("%s %v", q.Query(), q.args)
}

type Selection struct {
	Alias      string
	ColumnName string
	TableName  string
	TableAlias string
	StructName string
	FieldName  string
}

type DBContext interface {
	Placeholder(typ string, prevArgs []interface{}) string
	QuoteIdent(v string) string
}

type Querier interface {
	Apply(q *Query, ctx DBContext)
	Selection() Selection
}

type Aliased interface {
	Querier
	Alias() string
}

type Expr interface {
	Querier

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

	In(vals ...interface{}) PredExpr

	Add(v interface{}) AnonExpr
	Sbt(v interface{}) AnonExpr
	Mlt(v interface{}) AnonExpr
	Dvd(v interface{}) AnonExpr
	Concat(s interface{}) AnonExpr
}

// AnonExpr represents an anonymous (not aliased) expression.
type AnonExpr interface {
	Expr
	As(alias string) Aliased
}

// PredExpr represents this expression is a predicate.
type PredExpr interface {
	AnonExpr
	ImplPredExpr()
}

type TableLike interface {
	ApplyTable(q *Query, ctx DBContext)
}

type QueryTable interface {
	Querier
	TableLike
}

type SchemaTable interface {
	TableLike
	All() *ExprListExpr
	Columns() []*Column
}

type QueryExpr interface {
	Expr
	As(alias string) QueryTable
	Selections() []Selection
	Construct() Query
	OrderBy(ords ...Orderer) QueryExpr
	Limit(n int) QueryExpr
	Offset(n int) QueryExpr

	// Shallow copy QueryExpr and set LIMIT and OFFSET
	WithLimits(limit, offset int) QueryExpr
}

type SelectClause interface {
	QueryExpr
	From(table TableLike, tables ...TableLike) Clauses
}

type Clauses interface {
	QueryExpr
	Joins(joins ...JoinDefiner) Clauses
	Where(preds ...PredExpr) Clauses
	GroupBy(exps ...Expr) GroupByClause
}

type GroupByClause interface {
	QueryExpr
	Having(preds ...PredExpr) GroupByClause
}

type Orderer interface {
	Ordering() Ordering
}

const (
	ORDER_ASC  = "ASC"
	ORDER_DESC = "DESC"
)

type Order string

type Ordering struct {
	expr  Expr
	order Order
}

func (ord Ordering) Ordering() Ordering {
	return ord
}
