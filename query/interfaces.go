package query

type Built struct {
	Query string
	Args  []interface{}
}

type ExprBase interface {
	// String() string
	Build() Built
}

type Expr interface {
	ExprBase
	As(alias string) ExprAliased
}

type ExprAliased interface {
	ExprBase
	Alias() string // Name()?
	Expr() Expr
}

type PredExpr interface {
	Expr
	PredExpr()
}

// 集計関数など Predicate じゃない Expr を受け取りたい
// 箇所もあるから、Expr = Predicate | Operation でもいいかも

type QueryGlobal interface {
	Case() PredExpr // take special stuff
	Exists(q Query) PredExpr
	In(exps ...Expr) PredExpr
	And(exps ...PredExpr) PredExpr
	Or(exps ...PredExpr) PredExpr
	Not(exp PredExpr) PredExpr
	Func(name string, args ...Expr) FuncExpr

	Count(exp Expr) FuncExpr
	Sum(exp Expr) FuncExpr
	Min(exp Expr) FuncExpr
	Max(exp Expr) FuncExpr
	Avg(exp Expr) FuncExpr

	InnerJoin(table TableBase) Join
	LeftJoin(table TableBase) Join
	RightJoin(table TableBase) Join
	FullJoin(table TableBase) Join

	Raw(v interface{}) Expr // ?
	Str(s string) Expr      // surround by quotes
	// Date, DateTime?

	Parens(exp Expr) Expr
}

type ValExpr interface {
	Expr

	Eq(exp interface{}) PredExpr
	Neq(exp interface{}) PredExpr
	Gt(exp interface{}) PredExpr
	Gte(exp interface{}) PredExpr
	Lt(exp interface{}) PredExpr
	Lte(exp interface{}) PredExpr

	Like(s string) PredExpr
	Between(from interface{}, to interface{}) PredExpr
	IsNull() PredExpr
	IsNotNull() PredExpr

	Add(exp Expr) Expr
	Sbt(exp Expr) Expr
	Mlt(exp Expr) Expr
	Dvd(exp Expr) Expr
	Negate(exp Expr) Expr

	Concat(s Expr) Expr
}

type FuncExpr interface {
	Expr
	Name() string
}

type ColumnExpr interface {
	ValExpr
	ColumnName() string
	TableName() string
	StructName() string
	FieldName() string
}

type ColumnListExpr interface {
	Expr
	Columns() []ColumnExpr
}

type TableBase interface {
	Name() string
	All() ColumnListExpr
}

type Table interface {
	TableBase
	As(alias string) TableAliased
}

type TableAliased interface {
	TableBase
	Alias() string
}

type Query interface {
	Expr // Query 自体も式になる
	GetDetail() QueryData
}
type QueryData struct{}

type SelectClause interface {
	Select(exps ...ExprBase) FromClause
}

type FromClause interface {
	Query
	From(table TableBase, tables ...TableBase) Clauses
}

type Clauses interface {
	ExtraClauses

	// Where や Joins だけ独立して定義したいケースもあるかも
	Where(exps ...PredExpr) Clauses
	Joins(joins ...JoinOn) Clauses
	GroupBy(exps ...Expr) GroupQuery

	Map(f func(q Clauses) Clauses) Clauses
}

type Join interface {
	On(exp PredExpr) JoinOn
}
type JoinOn interface {
	Expr
	joinOn()
}

type GroupQuery interface {
	ExtraClauses
	Having(exps ...PredExpr) GroupQuery
}

type ExtraClauses interface {
	Query
	OrderBy(exps ...Expr) ExtraClauses
	Limit(n int) ExtraClauses
	Offset(n int) ExtraClauses
}
