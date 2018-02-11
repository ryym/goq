package cllct

import (
	"reflect"
)

type RowMapCollector struct {
	mp       *map[string]interface{}
	colNames []string
	targets  []int
}

func (cl *RowMapCollector) ImplSingleCollector() {}

func (cl *RowMapCollector) Init(conf *InitConf) (bool, error) {
	cl.colNames = conf.ColNames
	for i, _ := range conf.ColNames {
		if conf.take(i) {
			cl.targets = append(cl.targets, i)
		}
	}
	return true, nil
}

func (cl *RowMapCollector) AfterInit(conf *InitConf) error {
	return nil
}

func (cl *RowMapCollector) Next(ptrs []interface{}) {
	for _, i := range cl.targets {
		ptrs[i] = new(interface{})
	}
}

func (cl *RowMapCollector) AfterScan(ptrs []interface{}) {
	row := make(map[string]interface{}, len(ptrs))
	for _, i := range cl.targets {
		row[cl.colNames[i]] = reflect.ValueOf(ptrs[i]).Elem().Interface()
	}
	*cl.mp = row
}
