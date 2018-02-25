package goql

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type Query struct {
	query []string
	args  []interface{}
	errs  []error
}

func (q *Query) Query() string {
	return strings.Join(q.query, "")
}

func (q *Query) Args() []interface{} {
	return q.args
}

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

func (q Query) String() string {
	if len(q.errs) > 0 {
		return fmt.Sprintf("ERR: %s", q.Err())
	}
	return fmt.Sprintf("%s %v", q.Query(), q.args)
}
