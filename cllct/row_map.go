package cllct

import (
	"reflect"

	"github.com/ryym/goq/gql"
)

type RowMapCollector struct {
	mp       *map[string]interface{}
	colNames []string
}

func (c *RowMapCollector) ImplSingleCollector() {}

func (c *RowMapCollector) Init(_selects []gql.Selection, names []string) bool {
	c.colNames = names
	return true
}

func (c *RowMapCollector) Next(ptrs []interface{}) {
	for i := 0; i < len(ptrs); i++ {
		if ptrs[i] == nil {
			ptrs[i] = new(interface{})
		}
	}
}

func (c *RowMapCollector) AfterScan(ptrs []interface{}) {
	row := make(map[string]interface{}, len(ptrs))
	for i, p := range ptrs {
		row[c.colNames[i]] = reflect.ValueOf(p).Elem().Interface()
	}
	*c.mp = row
}
