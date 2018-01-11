package gql

type ops struct {
	expr Expr
}

func (o *ops) As(alias string) Aliased {
	return &aliased{o.expr, alias}
}
