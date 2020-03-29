package goq

import (
	"github.com/ryym/goq/dialect"
)

// Builder provides query builder methods and result collector methods.
type Builder struct {
	*QueryBuilder
	*CollectorMaker
}

func NewBuilder(dl dialect.Dialect) *Builder {
	return &Builder{
		QueryBuilder:   NewQueryBuilder(dl),
		CollectorMaker: NewMaker(),
	}
}
