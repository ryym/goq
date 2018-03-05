package cllct

import (
	"reflect"
)

// RowMapSliceCollector collects rows into a slice of maps.
//
//	[]map[string]interface{}{
//		{ "user_id": 30, "name": "alice" },
//		{ "user_id": 31, "name": "bob" },
//	}
//
// But be careful, this collector collects values without any conversions.
// See https://godoc.org/github.com/ryym/goq/cllct#RowMapCollector for details.
type RowMapSliceCollector struct {
	slice    *[]map[string]interface{}
	colNames []string
	targets  []int
}

func (cl *RowMapSliceCollector) ImplListCollector() {}

func (cl *RowMapSliceCollector) init(conf *initConf) (bool, error) {
	cl.colNames = conf.ColNames
	for i, _ := range conf.ColNames {
		if conf.take(i) {
			cl.targets = append(cl.targets, i)
		}
	}
	return true, nil
}

func (cl *RowMapSliceCollector) afterinit(conf *initConf) error {
	return nil
}

func (cl *RowMapSliceCollector) next(ptrs []interface{}) {
	for _, i := range cl.targets {
		ptrs[i] = new(interface{})
	}
}

func (cl *RowMapSliceCollector) afterScan(ptrs []interface{}) {
	row := make(map[string]interface{}, len(ptrs))
	for _, i := range cl.targets {
		row[cl.colNames[i]] = reflect.ValueOf(ptrs[i]).Elem().Interface()
	}
	*cl.slice = append(*cl.slice, row)
}
