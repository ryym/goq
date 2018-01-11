package gql

type Builder struct{}

func (b *Builder) Var(v interface{}) Expr {
	return (&litExpr{val: v}).init()
}
