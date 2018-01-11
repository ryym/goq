package gql

type Builder struct{}

func (b *Builder) Var(v interface{}) Expr {
	return (&litExpr{val: v}).init()
}

func (b *Builder) Raw(sql string) Expr {
	return (&rawExpr{sql: sql}).init()
}

func (b *Builder) Parens(exp Expr) Expr {
	return (&parensExpr{exp: exp}).init()
}
