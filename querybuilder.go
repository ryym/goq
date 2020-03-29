package goq

import "github.com/ryym/goq/dialect"

// QueryBuilder is a core query builder.
// It provides basic clauses and operators.
type QueryBuilder struct {
	ctx DBContext
}

func NewQueryBuilder(dl dialect.Dialect) *QueryBuilder {
	return &QueryBuilder{dl}
}

// Query constructs a Query from the given expression.
//
//	z := goql.NewQueryBuilder(dialect.Generic())
//	q := z.Query(z.Var(1).Add(z.Var(40)))
//	fmt.Println(q.String())
func (b *QueryBuilder) Query(exp QueryApplier) Query {
	q := Query{}
	exp.Apply(&q, b.ctx)
	return q
}

// Var is a variable for the query.
// This will be replaced by a placeholder and the value
// will be stored as an argument for the query.
func (b *QueryBuilder) Var(v interface{}) AnonExpr {
	return (&litExpr{val: v}).init()
}

// VarT is a variable with the type.
// PostgreSQL requires a type of placeholder in some situation.
func (b *QueryBuilder) VarT(v interface{}, typ string) AnonExpr {
	return (&litExpr{val: v, typ: typ}).init()
}

func (b *QueryBuilder) Null() AnonExpr {
	return (&nullExpr{}).init()
}

// Raw constructs a raw expression.
// The given string will be embedded to the query without any filterings.
func (b *QueryBuilder) Raw(sql string) *RawExpr {
	return (&RawExpr{sql: sql}).init()
}

// Parens surrounds the given expression by parentheses.
func (b *QueryBuilder) Parens(exp Expr) AnonExpr {
	return (&parensExpr{exp: exp}).init()
}

// Name is used to an alias for some expression.
// You can use this in 'ORDER BY' to use the alias of
// expression in 'SELECT' clause.
func (b *QueryBuilder) Name(name string) *nameExpr {
	return (&nameExpr{name: name}).init()
}

// Col creates a column expression dynamically.
func (b *QueryBuilder) Col(table, col string) *Column {
	return (&Column{
		tableName:  "",
		tableAlias: table,
		structName: "",
		name:       col,
		fieldName:  "",
	}).init()
}

// Table creates a table name dynamically.
func (b *QueryBuilder) Table(name string) *DynmTable {
	return newDynmTable(name)
}

// And concatenates all the given predicates by 'AND'.
func (b *QueryBuilder) And(preds ...PredExpr) PredExpr {
	return (&logicalOp{op: "AND", preds: preds}).init()
}

// Or concatenates all the given predicates by 'OR'.
func (b *QueryBuilder) Or(preds ...PredExpr) PredExpr {
	return (&logicalOp{op: "OR", preds: preds}).init()
}

// Func creates a function with its name and the arguments.
func (b *QueryBuilder) Func(name string, args ...interface{}) AnonExpr {
	expArgs := make([]Expr, 0, len(args))
	for _, a := range args {
		expArgs = append(expArgs, lift(a))
	}
	return (&funcExpr{name: name, args: expArgs}).init()
}

// Count counts the expression by 'COUNT'.
func (b *QueryBuilder) Count(exp Expr) AnonExpr {
	return b.Func("COUNT", exp)
}

// Sum gets a total sum of the expression by 'SUM'.
func (b *QueryBuilder) Sum(exp Expr) AnonExpr {
	return b.Func("SUM", exp)
}

// Min gets a minimum value of the expression by 'MIN'.
func (b *QueryBuilder) Min(exp Expr) AnonExpr {
	return b.Func("MIN", exp)
}

// Max gets a maximum value of the expression by 'MAX'.
func (b *QueryBuilder) Max(exp Expr) AnonExpr {
	return b.Func("MAX", exp)
}

// Avg gets an average value of the expression by 'AVG'.
func (b *QueryBuilder) Avg(exp Expr) AnonExpr {
	return b.Func("AVG", exp)
}

// Coalesce sets an alternative value for 'NULL' by 'COALESCE'.
func (b *QueryBuilder) Coalesce(exp Expr, alt interface{}) AnonExpr {
	return b.Func("COALESCE", exp, lift(alt))
}

// Concat is a 'CONCAT' function.
func (b *QueryBuilder) Concat(exps ...interface{}) AnonExpr {
	return b.Func("CONCAT", exps...)
}

// Exists constructs an 'EXISTS' predicate by the given query.
func (b *QueryBuilder) Exists(query QueryExpr) PredExpr {
	return &predExpr{(&existsExpr{query: query}).init()}
}

// NotExists constructs an 'NOT EXISTS' predicate by the given query.
func (b *QueryBuilder) NotExists(query QueryExpr) PredExpr {
	return &predExpr{(&existsExpr{query: query, not: true}).init()}
}

// Select constructs a 'SELECT' clause.
func (b *QueryBuilder) Select(exps ...Selectable) SelectClause {
	return (&queryExpr{exps: exps, ctx: b.ctx}).init()
}

// SelectDistinct constructs a 'SELECT' clause with 'DISTINCT'.
func (b *QueryBuilder) SelectDistinct(exps ...Selectable) SelectClause {
	return (&queryExpr{exps: exps, ctx: b.ctx, distinct: true}).init()
}

// InnerJoin constructs an 'INNER JOIN' clause.
func (b *QueryBuilder) InnerJoin(table TableLike) *JoinClause {
	return &JoinClause{JOIN_INNER, table}
}

// LeftJoin constructs an 'LEFT OUTER JOIN' clause.
func (b *QueryBuilder) LeftJoin(table TableLike) *JoinClause {
	return &JoinClause{JOIN_LEFT, table}
}

// RightJoin constructs an 'RIGHT OUTER JOIN' clause.
func (b *QueryBuilder) RightJoin(table TableLike) *JoinClause {
	return &JoinClause{JOIN_RIGHT, table}
}

// FullJoin constructs an 'FULL OUTER JOIN' clause.
func (b *QueryBuilder) FullJoin(table TableLike) *JoinClause {
	return &JoinClause{JOIN_FULL, table}
}

// Case constructs a 'CASE' expression.
//
//	z := goql.NewQueryBuilder(dialect.Generic())
//	age := z.Col("users", "age")
//	q := z.Case(
//		z.When(age.Lt(20)).Then("under20"),
//		z.When(age.Lt(40)).Then("under40"),
//	).Else("above40")
//
// SQL:
//
//	CASE
//	  WHEN users.age > 20 THEN 'under20'
//	  WHEN users.age > 40 THEN 'under40'
//	  ELSE 'above40'
//	END
func (b *QueryBuilder) Case(cases ...*WhenClause) *CaseExpr {
	return (&CaseExpr{cases: cases}).init()
}

// CaseOf constructs a 'CASE' expression for the value.
//
//	z := goql.NewQueryBuilder(dialect.Generic())
//	id := z.Col("users", "id")
//	q := z.CaseOf(id,
//		z.When(1).Then("one"),
//		z.When(2).Then("two"),
//	).Else("other")
//
// SQL:
//
//	CASE users.id
//	  WHEN 1 THEN 'one'
//	  WHEN 2 THEN 'two'
//	  ELSE 'other'
//	END
func (b *QueryBuilder) CaseOf(val Expr, cases ...*WhenClause) *CaseExpr {
	return (&CaseExpr{val: val, cases: cases}).init()
}

// When is a 'WHEN' clause for a 'CASE' expression.
func (b *QueryBuilder) When(when interface{}) *WhenClause {
	return &WhenClause{when: lift(when)}
}

// Modifier clauses

// InsertInto constructs an 'INSERT' statement.
// If the 'cols' are omitted,
// all columns of the table are set.
func (b *QueryBuilder) InsertInto(table SchemaTable, cols ...*Column) *InsertMaker {
	return &InsertMaker{
		table: table,
		cols:  cols,
		ctx:   b.ctx,
	}
}

// Update constructs an 'UPDATE' statement.
func (b *QueryBuilder) Update(table SchemaTable) *UpdateMaker {
	return &UpdateMaker{
		table: table,
		ctx:   b.ctx,
	}
}

// DeleteFrom constructs a 'DELETE' statement.
func (b *QueryBuilder) DeleteFrom(table SchemaTable) *Delete {
	return &Delete{
		table: table,
		ctx:   b.ctx,
	}
}
