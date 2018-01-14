package cllct

import "github.com/ryym/goq/gql"

type Collector interface {
	Init(selects []gql.Selection, colNames []string) (mappable bool)
	Next(ptrs []interface{})
	AfterScan(ptrs []interface{})
}
