// Package cllct provides various collectors.
// A collector collects results fetched from DB as sql.Rows
// into a specified format like a slice, map, etc.
package cllct

import "github.com/ryym/goq/goql"

type InitConf struct {
	Selects  []goql.Selection
	ColNames []string
	takens   map[int]bool
}

func NewInitConf(selects []goql.Selection, colNames []string) *InitConf {
	return &InitConf{selects, colNames, map[int]bool{}}
}

func (c *InitConf) take(colIdx int) bool {
	ok := c.takens[colIdx]
	if !ok {
		c.takens[colIdx] = true
	}
	return !ok
}

func (c *InitConf) canTake(colIdx int) bool {
	return !c.takens[colIdx]
}

type Collector interface {
	Init(conf *InitConf) (collectable bool, err error)
	AfterInit(conf *InitConf) error
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

type tableInfo struct {
	structName string
	tableAlias string
}
