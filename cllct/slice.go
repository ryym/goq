package cllct

import (
	"reflect"

	"github.com/ryym/goq/util"
)

type SliceCollector struct {
	elemType reflect.Type
	colToFld map[int]int
	ptr      interface{}
	slice    reflect.Value
	row      reflect.Value
}

func (cl *SliceCollector) ImplListCollector() {}

func (cl *SliceCollector) Init(conf *InitConf) (bool, error) {
	if err := checkPtrKind(cl.ptr, reflect.Slice); err != nil {
		return false, err
	}

	cl.slice = reflect.ValueOf(cl.ptr).Elem()
	cl.ptr = nil

	cl.elemType = cl.slice.Type().Elem()
	cl.slice.Set(reflect.MakeSlice(reflect.SliceOf(cl.elemType), 0, 0))

	targets := map[string]int{}
	for i := 0; i < cl.elemType.NumField(); i++ {
		f := cl.elemType.Field(i)
		if f.PkgPath == "" { // exported
			targets[util.FldToCol(f.Name)] = i
		}
	}

	cl.colToFld = map[int]int{}
	for iC, c := range conf.ColNames {
		if iF, ok := targets[c]; ok && conf.take(iC) {
			cl.colToFld[iC] = iF
		}
	}
	return len(cl.colToFld) > 0, nil
}

func (cl *SliceCollector) AfterInit(conf *InitConf) error {
	return nil
}

func (cl *SliceCollector) Next(ptrs []interface{}) {
	row := reflect.New(cl.elemType).Elem()
	cl.row = row.Addr()
	for c, f := range cl.colToFld {
		ptrs[c] = row.Field(f).Addr().Interface()
	}
}

func (cl *SliceCollector) AfterScan(_ptrs []interface{}) {
	cl.slice.Set(reflect.Append(cl.slice, cl.row.Elem()))
}
