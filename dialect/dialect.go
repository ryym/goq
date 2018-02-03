package dialect

import "fmt"

type Dialect interface {
	Placeholder(typ string, prevArgs []interface{}) string
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

func (dl *generic) Placeholder(typ string, prevArgs []interface{}) string {
	return "?"
}

func (dl *generic) QuoteIdent(v string) string {
	return v
}

type postgres struct{}

func Postgres() *postgres {
	return &postgres{}
}

func (dl *postgres) Placeholder(typ string, prevArgs []interface{}) string {
	ph := fmt.Sprintf("$%d", len(prevArgs)+1)
	if typ != "" {
		ph += "::" + typ
	}
	return ph
}

func (dl *postgres) QuoteIdent(v string) string {
	return fmt.Sprintf(`"%s"`, v)
}

type sqlite struct{}

func Sqlite() *sqlite {
	return &sqlite{}
}

func (dl *sqlite) Placeholder(typ string, prevArgs []interface{}) string {
	return "?"
}

func (dl *sqlite) QuoteIdent(v string) string {
	return fmt.Sprintf(`"%s"`, v)
}
