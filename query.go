package goq

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Query is a result query constructed by the query builder.
// You can use this to run a query using Go's sql.DB.
//
//	q, _ := query.Construct()
//	db, _ := sql.Open("postgres", conn)
//	db.Query(q.Query(), q.Args()...)
type Query struct {
	query []string
	args  []interface{}
	errs  []error
}

// Query returns a constructed query string.
// The query may contain placeholders.
func (q *Query) Query() string {
	return strings.Join(q.query, "")
}

// Args returns values for placeholders.
func (q *Query) Args() []interface{} {
	return q.args
}

// Err returns an error occurred during the query construction.
// The error contains one or more error messages joined by '|'.
func (q *Query) Err() error {
	if len(q.errs) == 0 {
		return nil
	}

	msgs := make([]string, 0, len(q.errs))
	for _, err := range q.errs {
		msgs = append(msgs, err.Error())
	}
	return errors.New(strings.Join(msgs, " | "))
}

// String converts this Query to a string.
// The string contains a query, an arguments slice, and an error.
func (q Query) String() string {
	if len(q.errs) > 0 {
		return fmt.Sprintf("ERR: %s", q.Err())
	}
	return fmt.Sprintf("%s %v", q.Query(), q.args)
}
