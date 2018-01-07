package query

type Query struct {
	Query string
	Args  []interface{}
}

// カラムでなくエイリアスもついてなければ全て空
type SelectItem struct {
	Alias      string
	ColumnName string
	TableName  string
	TableAlias string
	StructName string
	FieldName  string
}

type Queryable interface {
	ToQuery() Query
	ToSelectItem() SelectItem
}

type Expr interface {
	Queryable
	As(alias string) ExprAliased

	Eq(v interface{}) PredExpr
	// Neq(v interface{}) PredExpr
	// Gt(v interface{}) PredExpr
	// Gte(v interface{}) PredExpr
	// Lt(v interface{}) PredExpr
	// Lte(v interface{}) PredExpr
	// Like(s string) PredExpr
	// Between(from interface{}, to interface{}) PredExpr
	// IsNull() PredExpr
	// IsNotNull() PredExpr

	Add(v interface{}) Expr
	// Sbt(v interface{}) Expr
	Mlt(v interface{}) Expr
	// Dvd(v interface{}) Expr
	// Negate(v interface{}) Expr
	// Concat(s interface{}) Expr
}

type PredExpr interface {
	Expr
	PredExpr()
}

type ExprAliased interface {
	Queryable
	Alias() string
}

type QueryBuilder interface {
	Select(exps ...Queryable) SelectClause

	// Case() PredExpr // take special stuff
	Exists(q Query) PredExpr
	In(exps ...Expr) PredExpr
	And(exps ...PredExpr) PredExpr
	Or(exps ...PredExpr) PredExpr
	Not(exp PredExpr) PredExpr
	Func(name string, args ...Expr) Expr

	Val(v interface{}) LitExpr
	Raw(v interface{}) Expr
	Parens(exp Expr) Expr

	Count(exp Expr) Expr
	Sum(exp Expr) Expr
	Min(exp Expr) Expr
	Max(exp Expr) Expr
	Avg(exp Expr) Expr
	Coalesce(exp Expr, alt interface{}) Expr

	// サブクエリも受け取れるべき
	InnerJoin(table Table) Join
	LeftJoin(table Table) Join
	RightJoin(table Table) Join
	FullJoin(table Table) Join
}

type LitExpr interface {
	Queryable
	LitExpr()
}

// Select 時の情報は ToSelectItem として抽象されるから
// これは不要..?
// type ColumnExpr interface {
// 	Queryable
// 	ColumnName() string
// 	TableName() string
// 	StructName() string
// 	FieldName() string
// }

// table.All() で必要
type ExprListExpr interface {
	Queryable
	Queryables() []Queryable
}

type Table interface {
	TableName() string
	Alias() string
	All() ExprListExpr
	// Columns()
}

// As() の戻り値はテーブルごとに違うから interface にはできない
// type Table interface {
// 	Table
// 	As(alias string) TableAliased
// }

// type TableAliased interface {
// 	Table
// 	Alias() string
// }

type QueryStmt interface {
	Queryable
	GetSelects() []SelectItem
}

type SelectClause interface {
	QueryStmt
	From(table Table, tables ...Table) Clauses
}

type ExtraClauses interface {
	QueryStmt
	OrderBy(exps ...Queryable) ExtraClauses
	Limit(n int) ExtraClauses
	Offset(n int) ExtraClauses
}

type Clauses interface {
	ExtraClauses

	// Pred 限定だと Raw で直接条件を書けない
	Where(exps ...PredExpr) Clauses
	Joins(joins ...JoinOn) Clauses
	GroupBy(exps ...Queryable) GroupQuery
}

type GroupQuery interface {
	ExtraClauses
	Having(exps ...PredExpr) GroupQuery
}

type Join interface {
	On(exp PredExpr) JoinOn
	// Using
}

const (
	JOIN_INNER = iota
	JOIN_LEFT
	JOIN_RIGHT
	JOIN_FULL
)

type JoinKind int

type JoinOn struct {
	On   PredExpr
	Kind JoinKind
}
