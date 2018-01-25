package cllct

import (
	"reflect"

	"github.com/ryym/goq/gql"
)

type RowMapSliceCollector struct {
	slice    *[]map[string]interface{}
	colNames []string
}

func (c *RowMapSliceCollector) ImplListCollector() {}

func (c *RowMapSliceCollector) Init(_selects []gql.Selection, names []string) bool {
	c.colNames = names
	return true
}

func (c *RowMapSliceCollector) Next(ptrs []interface{}) {
	for i := 0; i < len(ptrs); i++ {
		if ptrs[i] == nil {
			ptrs[i] = new(interface{})
		}
	}
}

func (c *RowMapSliceCollector) AfterScan(ptrs []interface{}) {
	row := make(map[string]interface{}, len(ptrs))
	for i, p := range ptrs {
		row[c.colNames[i]] = reflect.ValueOf(p).Elem().Interface()
	}
	*c.slice = append(*c.slice, row)
}
