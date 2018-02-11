package cllct

import (
	"reflect"

	"github.com/ryym/goq/gql"
)

type ModelSliceCollector struct {
	elemType reflect.Type
	cols     []*gql.Column
	table    tableInfo
	colToFld map[int]int
	ptr      interface{}
	slice    reflect.Value
	row      reflect.Value
}

func (cl *ModelSliceCollector) ImplListCollector() {}

func (cl *ModelSliceCollector) Init(conf *InitConf) (bool, error) {
	if err := checkPtrKind(cl.ptr, reflect.Slice); err != nil {
		return false, err
	}

	cl.slice = reflect.ValueOf(cl.ptr).Elem()
	cl.ptr = nil

	cl.colToFld = map[int]int{}
	for iC, c := range conf.Selects {
		if conf.canTake(iC) && isSameTable(c, cl.table) {
			for iF, f := range cl.cols {
				if c.FieldName == f.FieldName() {
					cl.colToFld[iC] = iF
					conf.take(iC)
				}
			}
		}
	}
	cl.elemType = cl.slice.Type().Elem()
	cl.slice.Set(reflect.MakeSlice(reflect.SliceOf(cl.elemType), 0, 0))
	return len(cl.colToFld) > 0, nil
}

func (cl *ModelSliceCollector) Next(ptrs []interface{}) {
	row := reflect.New(cl.elemType).Elem()
	cl.row = row.Addr()
	for c, f := range cl.colToFld {
		ptrs[c] = row.Field(f).Addr().Interface()
	}
}

func (cl *ModelSliceCollector) AfterScan(_ptrs []interface{}) {
	cl.slice.Set(reflect.Append(cl.slice, cl.row.Elem()))
}
