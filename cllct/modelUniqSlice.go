package cllct

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/ryym/goq/goql"
)

type ModelUniqSliceCollector struct {
	elemType    reflect.Type
	cols        []*goql.Column
	table       tableInfo
	colToFld    map[int]int
	ptr         interface{}
	slice       reflect.Value
	pkFieldName string
	pkIdx       int
	pks         map[interface{}]bool
	elem        *reflect.Value
}

func (cl *ModelUniqSliceCollector) ImplListCollector() {}

func (cl *ModelUniqSliceCollector) Init(conf *initConf) (bool, error) {
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

	if cl.pkFieldName == "" {
		return false, fmt.Errorf("primary key not defined for %s", cl.table.structName)
	}

	cl.pks = map[interface{}]bool{}

	// Since `ModelUniqSliceCollector` does not need to create a struct every row,
	// prepare only one struct passed to `Rows.Scan` as a pointer,
	// and copy this only if necessary.
	elem := reflect.New(cl.elemType).Elem()
	cl.elem = &elem

	cl.pkIdx = -1
	cl.colToFld = map[int]int{}
	for iC, c := range conf.Selects {
		if conf.canTake(iC) && isSameTable(c, cl.table) {
			if cl.pkFieldName == c.FieldName {
				cl.pkIdx = iC
			}
			for iF, f := range cl.cols {
				if c.FieldName == f.FieldName() {
					cl.colToFld[iC] = iF
					conf.take(iC)
				}
			}
		}
	}

	if cl.pkIdx == -1 {
		return false, fmt.Errorf("primary key %s not selected", cl.pkFieldName)
	}

	return len(cl.colToFld) > 0, nil
}

func (cl *ModelUniqSliceCollector) AfterInit(conf *initConf) error {
	return nil
}

func (cl *ModelUniqSliceCollector) Next(ptrs []interface{}) {
	for c, f := range cl.colToFld {
		ptrs[c] = cl.elem.Field(f).Addr().Interface()
	}
}

func (cl *ModelUniqSliceCollector) AfterScan(ptrs []interface{}) {
	pk := reflect.ValueOf(ptrs[cl.pkIdx]).Elem().Interface()
	if cl.pks[pk] {
		return
	}

	// Copy `elem` and append it.
	cl.slice.Set(reflect.Append(cl.slice, *cl.elem))
	cl.pks[pk] = true
}
