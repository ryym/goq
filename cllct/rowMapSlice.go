package cllct

import (
	"reflect"
)

type RowMapSliceCollector struct {
	slice    *[]map[string]interface{}
	colNames []string
	targets  []int
}

func (cl *RowMapSliceCollector) ImplListCollector() {}

func (cl *RowMapSliceCollector) Init(conf *initConf) (bool, error) {
	cl.colNames = conf.ColNames
	for i, _ := range conf.ColNames {
		if conf.take(i) {
			cl.targets = append(cl.targets, i)
		}
	}
	return true, nil
}

func (cl *RowMapSliceCollector) AfterInit(conf *initConf) error {
	return nil
}

func (cl *RowMapSliceCollector) Next(ptrs []interface{}) {
	for _, i := range cl.targets {
		ptrs[i] = new(interface{})
	}
}

func (cl *RowMapSliceCollector) AfterScan(ptrs []interface{}) {
	row := make(map[string]interface{}, len(ptrs))
	for _, i := range cl.targets {
		row[cl.colNames[i]] = reflect.ValueOf(ptrs[i]).Elem().Interface()
	}
	*cl.slice = append(*cl.slice, row)
}
