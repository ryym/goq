package dialect

import "fmt"

type Dialect interface {
	Placeholder(prevArgs []interface{}) string
	QuoteIdent(v string) string
}

func New(driver string) Dialect {
	switch driver {
	case "postgres":
		return &postgres{}
	case "sqlite3":
		return &sqlite{}
	}
	return nil
}

type generic struct{}

func Generic() *generic {
	return &generic{}
}

func (dl *generic) Placeholder(prevArgs []interface{}) string {
	return "?"
}

func (dl *generic) QuoteIdent(v string) string {
	return v
}

type postgres struct{}

func Postgres() *postgres {
	return &postgres{}
}

func (dl *postgres) Placeholder(prevArgs []interface{}) string {
	return fmt.Sprintf("$%d", len(prevArgs)+1)
}

func (dl *postgres) QuoteIdent(v string) string {
	return fmt.Sprintf(`"%s"`, v)
}

type sqlite struct{}

func Sqlite() *sqlite {
	return &sqlite{}
}

func (dl *sqlite) Placeholder(prevArgs []interface{}) string {
	return "?"
}

func (dl *sqlite) QuoteIdent(v string) string {
	return fmt.Sprintf(`"%s"`, v)
}
