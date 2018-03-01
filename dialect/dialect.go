// Package dialect defines some dialects per RDB.
// Currently it supports Postgres, MySQL, and SQLite3.
package dialect

import "fmt"

type Dialect interface {
	// Placeholder makes a placeholder string.
	Placeholder(typ string, prevArgs []interface{}) string

	// QuoteIndent surrounds an identifier by proper quotes.
	QuoteIdent(v string) string
}

// New returns a dialect for the given driver.
// It supports "postgres", "mysql", and "sqlite3".
// It returns nil if you pass an unknown driver name.
func New(driver string) Dialect {
	switch driver {
	case "postgres":
		return &postgres{}
	case "mysql":
		return &mysql{}
	case "sqlite3":
		return &sqlite{}
	}
	return nil
}

type generic struct{}

// Generic returns a generic dialect.
// Queries constructed by this dialect should work with most databases.
func Generic() Dialect {
	return &generic{}
}

func (dl *generic) Placeholder(typ string, prevArgs []interface{}) string {
	return "?"
}

func (dl *generic) QuoteIdent(v string) string {
	return v
}

type postgres struct{}

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

type mysql struct{}

func (dl *mysql) Placeholder(typ string, prevArgs []interface{}) string {
	return "?"
}

func (dl *mysql) QuoteIdent(v string) string {
	return fmt.Sprintf("`%s`", v)
}

type sqlite struct{}

func (dl *sqlite) Placeholder(typ string, prevArgs []interface{}) string {
	return "?"
}

func (dl *sqlite) QuoteIdent(v string) string {
	return fmt.Sprintf(`"%s"`, v)
}
