package goq

import (
	"errors"
	"fmt"
	"reflect"
)

// ModelUniqSliceCollector collects rows into a slice of models uniquely.
// The uniqueness of a model is determined by its primary key.
// If the result rows contains multiple rows which have a same primary key,
// ModelUniqSliceCollector scans the only first row.
//
// Example:
//
//	[]City{
//		{ ID: 8, Name: "Osaka" },
//		{ ID: 12, Name: "Kyoto" },
//	}
type ModelUniqSliceCollector struct {
	elemType    reflect.Type
	cols        []*Column
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

func (cl *ModelUniqSliceCollector) init(conf *initConf) (bool, error) {
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

func (cl *ModelUniqSliceCollector) afterinit(conf *initConf) error {
	return nil
}

func (cl *ModelUniqSliceCollector) next(ptrs []interface{}) {
	for c, f := range cl.colToFld {
		ptrs[c] = cl.elem.Field(f).Addr().Interface()
	}
}

func (cl *ModelUniqSliceCollector) afterScan(ptrs []interface{}) {
	pk := reflect.ValueOf(ptrs[cl.pkIdx]).Elem().Interface()
	if cl.pks[pk] {
		return
	}

	// Copy `elem` and append it.
	cl.slice.Set(reflect.Append(cl.slice, *cl.elem))
	cl.pks[pk] = true
}
