package dialect

import "fmt"

type Dialect interface {
	Placeholder(prevArgs []interface{}) string
	QuoteIdent(v string) string
}

func New(driver string) Dialect {
	switch driver {
	case "postgres":
		return &Postgres{}
	case "sqlite3":
		return &Sqlite{}
	}
	return nil
}

type Postgres struct{}

func (ctx *Postgres) Placeholder(prevArgs []interface{}) string {
	return fmt.Sprintf("$%d", len(prevArgs)+1)
}

func (ctx *Postgres) QuoteIdent(v string) string {
	return fmt.Sprintf(`"%s"`, v)
}

type Sqlite struct{}

func (ctx *Sqlite) Placeholder(prevArgs []interface{}) string {
	return "?"
}

func (ctx *Sqlite) QuoteIdent(v string) string {
	return fmt.Sprintf(`"%s"`, v)
}
