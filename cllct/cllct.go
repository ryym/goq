package cllct

import "github.com/ryym/goq/gql"

type InitConf struct {
	Selects  []gql.Selection
	ColNames []string
	takens   map[int]bool
}

func NewInitConf(selects []gql.Selection, colNames []string) *InitConf {
	return &InitConf{selects, colNames, map[int]bool{}}
}

func (c *InitConf) take(colIdx int) bool {
	ok := c.takens[colIdx]
	if !ok {
		c.takens[colIdx] = true
	}
	return !ok
}

type Collector interface {
	Init(conf *InitConf) (collectable bool, err error)
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
