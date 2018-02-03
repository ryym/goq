package cllct

import (
	"fmt"
	"reflect"

	"github.com/ryym/goq/gql"
)

type ModelUniqSliceCollector struct {
	elemType    reflect.Type
	cols        []*gql.Column
	structName  string
	tableAlias  string
	colToFld    map[int]int
	slice       reflect.Value
	pkFieldName string
	pkIdx       int
	pks         map[interface{}]bool
	elem        *reflect.Value
}

func (cl *ModelUniqSliceCollector) ImplListCollector() {}

func (cl *ModelUniqSliceCollector) Init(conf *InitConf) (bool, error) {
	if cl.pkFieldName == "" {
		return false, fmt.Errorf("primary key not defined for %s", cl.structName)
	}

	cl.pks = map[interface{}]bool{}
	cl.elemType = cl.slice.Type().Elem()
	cl.slice.Set(reflect.MakeSlice(reflect.SliceOf(cl.elemType), 0, 0))

	// Since `ModelUniqSliceCollector` does not need to create a struct every row,
	// prepare only one struct passed to `Rows.Scan` as a pointer,
	// and copy this only if necessary.
	elem := reflect.New(cl.elemType).Elem()
	cl.elem = &elem

	cl.pkIdx = -1
	cl.colToFld = map[int]int{}
	for iC, c := range conf.Selects {
		if conf.canTake(iC) && c.TableAlias == cl.tableAlias && c.StructName == cl.structName {
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
