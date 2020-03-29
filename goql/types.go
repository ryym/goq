package goql

// Selection provides information about a selected expression.
type Selection struct {
	Alias      string
	ColumnName string
	TableName  string
	TableAlias string
	StructName string
	FieldName  string
}

// DBContext abstructs query syntax differences among RDBs.
type DBContext interface {
	Placeholder(typ string, prevArgs []interface{}) string
	QuoteIdent(v string) string
}

// QueryApplier is the interface to append query parts.
// All query part structs implement this.
type QueryApplier interface {
	Apply(q *Query, ctx DBContext)
}

// Selectable represents a selectable expression.
// QueryBuilder.Select accepts values that implement this interface.
type Selectable interface {
	QueryApplier
	Selection() Selection
}

// Aliased is an aliased expression.
// You cannot operate an aliased expression.
// For example 'A.Add(B.As("b"))' fails to compile.
type Aliased interface {
	Selectable
	Alias() string
}

// Expr is the interface of expressions that provide basic operators.
// All expression structs implement this interface.
type Expr interface {
	Selectable

	// Eq does '=' comparison (equal to).
	Eq(v interface{}) PredExpr

	// Neq does '<>' comparison (not equal to).
	Neq(v interface{}) PredExpr

	// Gt does '>' comparison (greater than).
	Gt(v interface{}) PredExpr

	// Gte does '>=' comparison (greater than or equal to).
	Gte(v interface{}) PredExpr

	// Lt does '<' comparison (less than).
	Lt(v interface{}) PredExpr

	// Lte does '<=' comparison (less than or equal to).
	Lte(v interface{}) PredExpr

	// Add does '+' operation (addition).
	Add(v interface{}) AnonExpr

	// Sbt does '-' operation (subtraction).
	Sbt(v interface{}) AnonExpr

	// Mlt does '*' operation (multiplication).
	Mlt(v interface{}) AnonExpr

	// Dvd does '/' operation (division).
	Dvd(v interface{}) AnonExpr

	// Concat concats a string by '||'.
	Concat(s interface{}) AnonExpr

	// Like does 'LIKE' comparison.
	// For example, 'Title.Like("I am%")' becomes true if
	// the 'Title' starts with 'I am'.
	Like(s string) PredExpr

	// Between does 'BETWEEN' comparison.
	Between(start interface{}, end interface{}) PredExpr

	// In constructs 'IN' comparison.
	// You must pass a slice or query.
	// Otherwise a query construction will result in an error.
	In(valsOrQuery interface{}) PredExpr

	// IsNull does 'IS NULL' comparison.
	IsNull() PredExpr

	// IsNotNull does 'IS NOT NULL' comparison.
	IsNotNull() PredExpr
}

// AnonExpr represents an anonymous expression that does not have an alias.
type AnonExpr interface {
	Expr
	As(alias string) Aliased
}

// PredExpr represents a predicate expression.
// Some clauses like 'WHERE' or 'ON' accept predicate expressions only.
type PredExpr interface {
	AnonExpr
	ImplPredExpr()
}

// TableLike represents an expression that can be a table.
// Addition to database tables,
// aliased queries can also be used for 'FROM' or 'JOIN'.
// This means a query implements TableLike as well.
type TableLike interface {
	ApplyTable(q *Query, ctx DBContext)
}

// QueryTable represents a query as a table.
type QueryTable interface {
	Selectable
	TableLike
}

// SchemaTable represents a database table.
type SchemaTable interface {
	TableLike

	// All constructs the column list expression that has all columns.
	All() *ColumnList

	// Except constructs the column list expression that has all columns except columns you specify.
	Except(cols ...*Column) *ColumnList
}

// QueryRoot is the interface of query statements
// such as 'SELECT', 'CREATE', 'UPDATE', and so on.
type QueryRoot interface {
	// Construct constructs a query.
	Construct() (Query, error)
}

// QueryExpr is the interface of a query to select data.
type QueryExpr interface {
	Expr
	QueryRoot

	// Select overrides the selections.
	Select(exps ...Selectable) Clauses

	// As gives an alias to the query.
	As(alias string) QueryTable

	// Selections lists selected data.
	Selections() []Selection

	// OrderBy adds 'ORDER BY' clause.
	OrderBy(ords ...Orderer) QueryExpr

	// Limit adds 'LIMIT' clause.
	Limit(n int) QueryExpr

	// Offset adds 'OFFSET' clause.
	Offset(n int) QueryExpr

	// WithLimits copy this QueryExpr shallowly and set 'LIMIT' and 'OFFSET'.
	WithLimits(limit, offset int) QueryExpr
}

// SelectClause is a 'SELECT'  clause.
type SelectClause interface {
	QueryExpr

	// From constructs a 'FROM' clause.
	From(table TableLike, tables ...TableLike) Clauses
}

// Clauses has 'JOIN', 'WHERE', and 'GROUP BY' clauses.
// You can call 'WHERE' and 'JOIN' multiple times.
// All given expressions are appended to the query.
type Clauses interface {
	QueryExpr

	// Joins constructs a 'JOIN' clauses.
	Joins(joins ...JoinDefiner) Clauses

	// Where constructs a 'WHERE' clause.
	Where(preds ...PredExpr) Clauses

	// GroupBy constructs a 'GROUP BY' clasue.
	GroupBy(exps ...Expr) GroupByClause
}

// GroupByClause is a 'GROUP BY' clause.
type GroupByClause interface {
	QueryExpr

	// Having constructs a 'HAVING' clause.
	Having(preds ...PredExpr) GroupByClause
}

// Orderer is the interface of expressions for 'ORDER BY'.
type Orderer interface {
	Ordering() Ordering
}

// Orders.
const (
	ORDER_ASC  = "ASC"
	ORDER_DESC = "DESC"
)

type Order string

// Ordering is a combination of an expression and an order.
type Ordering struct {
	expr  Expr
	order Order
}

func (ord Ordering) Ordering() Ordering {
	return ord
}
