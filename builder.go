package goq

import (
	"github.com/ryym/goq/cllct"
	"github.com/ryym/goq/gql"
)

type Builder struct {
	*gql.Builder
	*cllct.CollectorMaker
}
