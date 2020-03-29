package goq

import (
	"reflect"

	"github.com/ryym/goq/util"
)

// SliceCollector collects rows into a slice of structs.
//
// Example:
//
//	[]IDs{
//		{ CountryID: 1, CityID: 3, AddressID: 18 },
//		{ CountryID: 1, CityID: 5, AddressID: 224 },
//	}
type SliceCollector struct {
	elemType reflect.Type
	colToFld map[int]int
	ptr      interface{}
	slice    reflect.Value
	row      reflect.Value
}

func (cl *SliceCollector) ImplListCollector() {}

func (cl *SliceCollector) init(conf *initConf) (bool, error) {
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

func (cl *SliceCollector) afterinit(conf *initConf) error {
	return nil
}

func (cl *SliceCollector) next(ptrs []interface{}) {
	row := reflect.New(cl.elemType).Elem()
	cl.row = row.Addr()
	for c, f := range cl.colToFld {
		ptrs[c] = row.Field(f).Addr().Interface()
	}
}

func (cl *SliceCollector) afterScan(_ptrs []interface{}) {
	cl.slice.Set(reflect.Append(cl.slice, cl.row.Elem()))
}
