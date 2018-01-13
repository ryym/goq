package gql

import (
	"fmt"
	"strings"
)

type Query struct {
	query []string
	args  []interface{}
}

func (q Query) String() string {
	return fmt.Sprintf("%s %v", strings.Join(q.query, ""), q.args)
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
	Placeholder(prevArgs []interface{}) string
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

type ExprListExpr interface {
	Querier
	Exprs() []Expr
}

type Table interface {
	TableName() string
	TableAlias() string
	All() ExprListExpr
	Columns() []Column
}

type Column interface {
	Expr
	TableName() string
	TableAlias() string
	StructName() string
	ColumnName() string
	FieldName() string
}

type QueryExpr interface {
	Expr
	Selections() []Selection
}

type SelectClause interface {
	QueryExpr
	From(table Table, tables ...Table) Clauses
}

type Clauses interface {
	QueryExpr
	Joins(joins ...JoinOn) Clauses
	Where(preds ...PredExpr) Clauses
	GroupBy(exps ...Expr) GroupByClause
}

type GroupByClause interface {
	QueryExpr
}

type JoinClause struct {
	joinType JoinType
	table    Table
}

func (jc *JoinClause) On(pred PredExpr) JoinOn {
	return JoinOn{jc.table, pred, jc.joinType}
}

const (
	JOIN_INNER = "INNER"
	JOIN_LEFT  = "LEFT OUTER"
	JOIN_RIGHT = "RIGHT OUTER"
	JOIN_FULL  = "FULL OUTER"
)

type JoinType string

type JoinOn struct {
	Table Table
	On    PredExpr
	Type  JoinType
}
