package goq

import (
	"reflect"
)

// RowMapCollector scans a first row into a map.
//
//	map[string]interface{ "user_id": 30, "name": "alice" }
//
// The keys of map are column name and values are row vlaues.
// So you can scan a row without defining a struct for it.
// But be careful, this scans values wihtout any conversions.
// Also result value types may be different by RDB.
// For example, if you scan an integer value '12' named 'id' from MySQL,
// the result map will be '{ "id": []uint8{49, 50} }', not '{ "id": 12 }'.
// This is because the MySQL driver returns '12' as bytes of UTF-8.
// However, the PostgreSQL driver returns '12' as an int64 value.
type RowMapCollector struct {
	mp       *map[string]interface{}
	colNames []string
	targets  []int
}

func (cl *RowMapCollector) ImplSingleCollector() {}

func (cl *RowMapCollector) init(conf *initConf) (bool, error) {
	cl.colNames = conf.ColNames
	for i, _ := range conf.ColNames {
		if conf.take(i) {
			cl.targets = append(cl.targets, i)
		}
	}
	return true, nil
}

func (cl *RowMapCollector) afterinit(conf *initConf) error {
	return nil
}

func (cl *RowMapCollector) next(ptrs []interface{}) {
	for _, i := range cl.targets {
		ptrs[i] = new(interface{})
	}
}

func (cl *RowMapCollector) afterScan(ptrs []interface{}) {
	row := make(map[string]interface{}, len(ptrs))
	for _, i := range cl.targets {
		row[cl.colNames[i]] = reflect.ValueOf(ptrs[i]).Elem().Interface()
	}
	*cl.mp = row
}
