package cllct

import (
	"reflect"

	"github.com/ryym/goq/goql"
)

// ElemCollector scans a first row into a model struct.
type ModelElemCollector struct {
	cols     []*goql.Column
	table    tableInfo
	colToFld map[int]int
	ptr      interface{}
	elem     reflect.Value
}

func (cl *ModelElemCollector) ImplSingleCollector() {}

func (cl *ModelElemCollector) init(conf *initConf) (bool, error) {
	if err := checkPtrKind(cl.ptr, reflect.Struct); err != nil {
		return false, err
	}

	cl.elem = reflect.ValueOf(cl.ptr).Elem()
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
	return len(cl.colToFld) > 0, nil
}

func (cl *ModelElemCollector) afterinit(conf *initConf) error {
	return nil
}

func (cl *ModelElemCollector) next(ptrs []interface{}) {
	for c, f := range cl.colToFld {
		ptrs[c] = cl.elem.Field(f).Addr().Interface()
	}
}

func (cl *ModelElemCollector) afterScan(_ptrs []interface{}) {}
