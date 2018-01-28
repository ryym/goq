package cllct

import (
	"reflect"
)

type RowMapSliceCollector struct {
	slice    *[]map[string]interface{}
	colNames []string
}

func (c *RowMapSliceCollector) ImplListCollector() {}

func (c *RowMapSliceCollector) Init(conf *InitConf) (bool, error) {
	for i, col := range conf.ColNames {
		if conf.take(i) {
			c.colNames = append(c.colNames, col)
		}
	}
	return true, nil
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
