package goql

import "github.com/ryym/goq/dialect"

// Builder is a core query builder.
// It provides basic clauses and operators.
type Builder struct {
	ctx DBContext
}

func NewBuilder(dl dialect.Dialect) *Builder {
	return &Builder{dl}
}

// Query constructs a Query from the given expression.
//
//	z := goql.NewBuilder(dialect.Generic())
//	q := z.Query(z.Var(1).Add(z.Var(40)))
//	fmt.Println(q.String())
func (b *Builder) Query(exp QueryApplier) Query {
	q := Query{}
	exp.Apply(&q, b.ctx)
	return q
}

// Var is a variable for the query.
// This will be replaced by a placeholder and the value
// will be stored as an argument for the query.
func (b *Builder) Var(v interface{}) AnonExpr {
	return (&litExpr{val: v}).init()
}

// VarT is a variable with the type.
// PostgreSQL requires a type of placeholder in some situation.
func (b *Builder) VarT(v interface{}, typ string) AnonExpr {
	return (&litExpr{val: v, typ: typ}).init()
}

func (b *Builder) Null() AnonExpr {
	return (&nullExpr{}).init()
}

// Raw constructs a raw expression.
// The given string will be embedded to the query without any filterings.
func (b *Builder) Raw(sql string) AnonExpr {
	return (&rawExpr{sql: sql}).init()
}

// Parens surrounds the given expression by parentheses.
func (b *Builder) Parens(exp Expr) AnonExpr {
	return (&parensExpr{exp: exp}).init()
}

// Name is used to an alias for some expression.
// You can use this in 'ORDER BY' to use the alias of
// expression in 'SELECT' clause.
func (b *Builder) Name(name string) *nameExpr {
	return (&nameExpr{name: name}).init()
}

// Col creates a column expression dynamically.
func (b *Builder) Col(table, col string) *Column {
	return (&Column{
		tableName:  "",
		tableAlias: table,
		structName: "",
		name:       col,
		fieldName:  "",
	}).init()
}

// Table creates a table name dynamically.
func (b *Builder) Table(name string) *DynmTable {
	return newDynmTable(name)
}

// And concatenates all the given predicates by 'AND'.
func (b *Builder) And(preds ...PredExpr) PredExpr {
	return (&logicalOp{op: "AND", preds: preds}).init()
}

// Or concatenates all the given predicates by 'OR'.
func (b *Builder) Or(preds ...PredExpr) PredExpr {
	return (&logicalOp{op: "OR", preds: preds}).init()
}

// Func creates a function with its name and the arguments.
func (b *Builder) Func(name string, args ...interface{}) AnonExpr {
	expArgs := make([]Expr, 0, len(args))
	for _, a := range args {
		expArgs = append(expArgs, lift(a))
	}
	return (&funcExpr{name: name, args: expArgs}).init()
}

// Count counts the expression by 'COUNT'.
func (b *Builder) Count(exp Expr) AnonExpr {
	return b.Func("COUNT", exp)
}

// Sum gets a total sum of the expression by 'SUM'.
func (b *Builder) Sum(exp Expr) AnonExpr {
	return b.Func("SUM", exp)
}

// Min gets a minimum value of the expression by 'MIN'.
func (b *Builder) Min(exp Expr) AnonExpr {
	return b.Func("MIN", exp)
}

// Max gets a maximum value of the expression by 'MAX'.
func (b *Builder) Max(exp Expr) AnonExpr {
	return b.Func("MAX", exp)
}

// Avg gets an average value of the expression by 'AVG'.
func (b *Builder) Avg(exp Expr) AnonExpr {
	return b.Func("AVG", exp)
}

// Coalesce sets an alternative value for 'NULL' by 'COALESCE'.
func (b *Builder) Coalesce(exp Expr, alt interface{}) AnonExpr {
	return b.Func("COALESCE", exp, lift(alt))
}

// Concat is a 'CONCAT' function.
func (b *Builder) Concat(exps ...interface{}) AnonExpr {
	return b.Func("CONCAT", exps...)
}

// Exists constructs an 'EXISTS' predicate by the given query.
func (b *Builder) Exists(query QueryExpr) PredExpr {
	return &predExpr{(&existsExpr{query: query}).init()}
}

// NotExists constructs an 'NOT EXISTS' predicate by the given query.
func (b *Builder) NotExists(query QueryExpr) PredExpr {
	return &predExpr{(&existsExpr{query: query, not: true}).init()}
}

// Select constructs a 'SELECT' clause.
func (b *Builder) Select(exps ...Selectable) SelectClause {
	return (&queryExpr{exps: exps, ctx: b.ctx}).init()
}

// InnerJoin constructs an 'INNER JOIN' clause.
func (b *Builder) InnerJoin(table TableLike) *JoinClause {
	return &JoinClause{JOIN_INNER, table}
}

// LeftJoin constructs an 'LEFT OUTER JOIN' clause.
func (b *Builder) LeftJoin(table TableLike) *JoinClause {
	return &JoinClause{JOIN_LEFT, table}
}

// RightJoin constructs an 'RIGHT OUTER JOIN' clause.
func (b *Builder) RightJoin(table TableLike) *JoinClause {
	return &JoinClause{JOIN_RIGHT, table}
}

// FullJoin constructs an 'FULL OUTER JOIN' clause.
func (b *Builder) FullJoin(table TableLike) *JoinClause {
	return &JoinClause{JOIN_FULL, table}
}

// Case constructs a 'CASE' expression.
//
//	z := goql.NewBuilder(dialect.Generic())
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
func (b *Builder) Case(cases ...*WhenClause) *CaseExpr {
	return (&CaseExpr{cases: cases}).init()
}

// CaseOf constructs a 'CASE' expression for the value.
//
//	z := goql.NewBuilder(dialect.Generic())
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
func (b *Builder) CaseOf(val Expr, cases ...*WhenClause) *CaseExpr {
	return (&CaseExpr{val: val, cases: cases}).init()
}

// When is a 'WHEN' clause for a 'CASE' expression.
func (b *Builder) When(when interface{}) *WhenClause {
	return &WhenClause{when: lift(when)}
}

// Modifier clauses

// InsertInto constructs an 'INSERT' statement.
// If the 'cols' are omitted,
// all columns of the table are set.
func (b *Builder) InsertInto(table SchemaTable, cols ...*Column) *InsertMaker {
	return &InsertMaker{
		table: table,
		cols:  cols,
		ctx:   b.ctx,
	}
}

// Update constructs an 'UPDATE' statement.
func (b *Builder) Update(table SchemaTable) *UpdateMaker {
	return &UpdateMaker{
		table: table,
		ctx:   b.ctx,
	}
}

// DeleteFrom constructs a 'DELETE' statement.
func (b *Builder) DeleteFrom(table SchemaTable) *Delete {
	return &Delete{
		table: table,
		ctx:   b.ctx,
	}
}
