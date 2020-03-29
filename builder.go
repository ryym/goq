package goq

import (
	"github.com/ryym/goq/dialect"
	"github.com/ryym/goq/goql"
)

// Builder provides query builder methods and result collector methods.
type Builder struct {
	*goql.QueryBuilder
	*CollectorMaker
}

func NewBuilder(dl dialect.Dialect) *Builder {
	return &Builder{
		QueryBuilder:   goql.NewQueryBuilder(dl),
		CollectorMaker: NewMaker(),
	}
}
