package goq

import "errors"

var (
	ErrNoRows = errors.New("goq: no rows in result set")
)
