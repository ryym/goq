package goq

import (
	"errors"
	"reflect"

	"github.com/ryym/goq/goql"
)

// ModelSliceCollector collects rows into a slice of models.
//
// Example:
//
//	[]City{
//		{ ID: 8, Name: "Osaka" },
//		{ ID: 12, Name: "Kyoto" },
//	}
type ModelSliceCollector struct {
	elemType reflect.Type
	cols     []*goql.Column
	table    tableInfo
	colToFld map[int]int
	ptr      interface{}
	slice    reflect.Value
	row      reflect.Value
}

func (cl *ModelSliceCollector) ImplListCollector() {}

func (cl *ModelSliceCollector) init(conf *initConf) (bool, error) {
	if err := checkPtrKind(cl.ptr, reflect.Slice); err != nil {
		return false, err
	}

	cl.slice = reflect.ValueOf(cl.ptr).Elem()
	cl.ptr = nil

	cl.elemType = cl.slice.Type().Elem()
	if cl.elemType.Kind() != reflect.Struct {
		return false, errors.New("slice elem type must be struct")
	}
	cl.slice.Set(reflect.MakeSlice(reflect.SliceOf(cl.elemType), 0, 0))

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
	return len(cl.colToFld) > 0, nil
}

func (cl *ModelSliceCollector) afterinit(conf *initConf) error {
	return nil
}

func (cl *ModelSliceCollector) next(ptrs []interface{}) {
	row := reflect.New(cl.elemType).Elem()
	cl.row = row.Addr()
	for c, f := range cl.colToFld {
		ptrs[c] = row.Field(f).Addr().Interface()
	}
}

func (cl *ModelSliceCollector) afterScan(_ptrs []interface{}) {
	cl.slice.Set(reflect.Append(cl.slice, cl.row.Elem()))
}
