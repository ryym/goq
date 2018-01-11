package gql

import "fmt"

type aliased struct {
	expr  Expr
	alias string
}

func (a *aliased) Alias() string { return a.alias }

func (a *aliased) Query() Query {
	qr := a.expr.Query()
	return Query{
		fmt.Sprintf("%s AS %s", qr.Query, a.alias),
		qr.Args,
	}
}

func (a *aliased) Selection() Selection { return Selection{} }
