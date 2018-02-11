package cllct

import (
	"errors"
	"reflect"

	"github.com/ryym/goq/gql"
)

type ModelMapCollector struct {
	elemType reflect.Type
	cols     []*gql.Column
	table    tableInfo
	colToFld map[int]int
	keySel   *gql.Selection
	keyIdx   int
	ptr      interface{}
	mp       reflect.Value
	row      reflect.Value
}

func (cl *ModelMapCollector) ImplListCollector() {}

func (cl *ModelMapCollector) Init(conf *InitConf) (bool, error) {
	if err := checkPtrKind(cl.ptr, reflect.Map); err != nil {
		return false, err
	}

	cl.mp = reflect.ValueOf(cl.ptr).Elem()
	cl.ptr = nil

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

	mapType := cl.mp.Type()

	cl.elemType = mapType.Elem()
	cl.mp.Set(reflect.MakeMap(reflect.MapOf(mapType.Key(), cl.elemType)))

	return len(cl.colToFld) > 0, nil
}

func (cl *ModelMapCollector) AfterInit(conf *InitConf) error {
	return nil
}

func (cl *ModelMapCollector) Next(ptrs []interface{}) {
	row := reflect.New(cl.elemType).Elem()
	cl.row = row.Addr()
	for c, f := range cl.colToFld {
		ptrs[c] = row.Field(f).Addr().Interface()
	}
}

func (cl *ModelMapCollector) AfterScan(ptrs []interface{}) {
	key := reflect.ValueOf(ptrs[cl.keyIdx]).Elem()
	cl.mp.SetMapIndex(key, cl.row.Elem())
}
