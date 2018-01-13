package gql

import "fmt"

type postgresCtx struct{}

func (ctx *postgresCtx) Placeholder(prevArgs []interface{}) string {
	return fmt.Sprintf("$%d", len(prevArgs)+1)
}

func (ctx *postgresCtx) QuoteIdent(v string) string {
	return fmt.Sprintf(`"%s"`, v)
}
