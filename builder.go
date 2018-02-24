package goq

import (
	"github.com/ryym/goq/cllct"
	"github.com/ryym/goq/dialect"
	"github.com/ryym/goq/goql"
)

type Builder struct {
	*goql.Builder
	*cllct.CollectorMaker
}

func NewBuilder(dl dialect.Dialect) *Builder {
	return &Builder{
		Builder:        goql.NewBuilder(dl),
		CollectorMaker: cllct.NewMaker(),
	}
}
