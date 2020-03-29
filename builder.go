package goq

import (
	"github.com/ryym/goq/dialect"
	"github.com/ryym/goq/goql"
)

// Builder provides query builder methods and result collector methods.
type Builder struct {
	*goql.Builder
	*CollectorMaker
}

func NewBuilder(dl dialect.Dialect) *Builder {
	return &Builder{
		Builder:        goql.NewBuilder(dl),
		CollectorMaker: NewMaker(),
	}
}
