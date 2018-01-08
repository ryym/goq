package gnr2

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

type Querier interface {
	Query() Query
	SelectItem() SelectItem
}

type QueryBuilder interface {
	Select(exps ...Querier) SelectClause

	// Case() PredExpr // take special stuff
	Exists(q Query) PredExpr
	In(exps ...Expr) PredExpr
	And(exps ...PredExpr) PredExpr
	Or(exps ...PredExpr) PredExpr
	Not(exp PredExpr) PredExpr
	Func(name string, args ...Expr) Expr

	Val(v interface{}) LitExpr
	Raw(v interface{}) RawExpr
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

type Expr interface {
	Querier
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
	// Mlt(v interface{}) Expr
	// Dvd(v interface{}) Expr
	// Negate(v interface{}) Expr
	// Concat(s interface{}) Expr
}

type PredExpr interface {
	Expr
	PredExpr()
}

type RawExpr interface {
	Expr
	RawExpr()
	Pred() PredExpr
}

type ExprAliased interface {
	Querier
	Alias() string
}

type LitExpr interface {
	Expr
	LitExpr()
}

// table.All() で必要
type ExprListExpr interface {
	Querier
	Exprs() []Expr
}

type QueryStmt interface {
	Querier
	GetSelects() []SelectItem
}

type SelectClause interface {
	QueryStmt
	From(table Table, tables ...Table) Clauses
}

type ExtraClauses interface {
	QueryStmt
	OrderBy(exps ...Expr) ExtraClauses
	Limit(n int) ExtraClauses
	Offset(n int) ExtraClauses
}

type Clauses interface {
	ExtraClauses
	Where(exps ...PredExpr) Clauses
	Joins(joins ...JoinOn) Clauses
	GroupBy(exps ...Expr) GroupQuery
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
	JOIN_INNER = "INNER"
	JOIN_LEFT  = "LEFT OUTER"
	JOIN_RIGHT = "RIGHT OUTER"
	JOIN_FULL  = "FULL OUTER"
)

type JoinKind string

type JoinOn struct {
	Table Table
	On    PredExpr
	Kind  JoinKind
}

type Table interface {
	TableName() string
	TableAlias() string
	All() ExprListExpr
	Columns() []Column
}

// テーブル struct の各カラムの interface。
// 結果マッピング時、SELECT された各式とのマッピングに必要。
type Column interface {
	Expr
	TableName() string
	TableAlias() string
	StructName() string
	ColumnName() string
	FieldName() string
}

type Collector interface {
	Init(selects []SelectItem, colNames []string) (mappable bool)
	Next(ptrs []interface{})
	AfterScan(ptrs []interface{})
}
