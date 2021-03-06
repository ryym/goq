package goq

import (
	"errors"
	"reflect"
)

// ModelMapCollector collects rows into a map of models
// whose key is a primary key.
type ModelMapCollector struct {
	elemType reflect.Type
	cols     []*Column
	table    tableInfo
	colToFld map[int]int
	keySel   *Selection
	keyIdx   int
	ptr      interface{}
	mp       reflect.Value
	row      reflect.Value
}

func (cl *ModelMapCollector) ImplListCollector() {}

func (cl *ModelMapCollector) init(conf *initConf) (bool, error) {
	if err := checkPtrKind(cl.ptr, reflect.Map); err != nil {
		return false, err
	}

	cl.mp = reflect.ValueOf(cl.ptr).Elem()
	cl.ptr = nil

	mapType := cl.mp.Type()
	cl.elemType = mapType.Elem()
	if cl.elemType.Kind() != reflect.Struct {
		return false, errors.New("map elem type must be struct")
	}

	cl.mp.Set(reflect.MakeMap(reflect.MapOf(mapType.Key(), cl.elemType)))

	if cl.keySel == nil {
		return false, errors.New("PK column required")
	}

	cl.colToFld = map[int]int{}
	cl.keyIdx = -1

	for iC, c := range conf.Selects {
		if conf.canTake(iC) && isSameTable(c, cl.table) {
			for iF, f := range cl.cols {
				if c.FieldName == f.FieldName() {
					cl.colToFld[iC] = iF
					conf.take(iC)
				}
			}
		}

		if isKeyCol(&c, cl.keySel) {
			cl.keyIdx = iC
		}
	}

	if cl.keyIdx == -1 {
		return false, errors.New("key not found")
	}
	if conf.canTake(cl.keyIdx) {
		return false, errors.New("PK column must be collected")
	}

	return len(cl.colToFld) > 0, nil
}

func (cl *ModelMapCollector) afterinit(conf *initConf) error {
	return nil
}

func (cl *ModelMapCollector) next(ptrs []interface{}) {
	row := reflect.New(cl.elemType).Elem()
	cl.row = row.Addr()
	for c, f := range cl.colToFld {
		ptrs[c] = row.Field(f).Addr().Interface()
	}
}

func (cl *ModelMapCollector) afterScan(ptrs []interface{}) {
	key := reflect.ValueOf(ptrs[cl.keyIdx]).Elem()
	cl.mp.SetMapIndex(key, cl.row.Elem())
}
