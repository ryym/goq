package cllct

import "github.com/ryym/goq/gql"

type Collector interface {
	Init(selects []gql.Selection, colNames []string) (collectable bool, err error)
	Next(ptrs []interface{})
	AfterScan(ptrs []interface{})
}

type ListCollector interface {
	Collector
	ImplListCollector()
}

type SingleCollector interface {
	Collector
	ImplSingleCollector()
}
