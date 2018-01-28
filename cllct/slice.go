package cllct

import (
	"reflect"

	"github.com/ryym/goq/gql"
	"github.com/ryym/goq/util"
)

type SliceCollector struct {
	elemType reflect.Type
	colToFld map[int]int
	slice    reflect.Value
	row      reflect.Value
}

func (cl *SliceCollector) ImplListCollector() {}

func (cl *SliceCollector) Init(selects []gql.Selection, names []string) (bool, error) {
	cl.elemType = cl.slice.Type().Elem()

	targets := map[string]int{}
	for i := 0; i < cl.elemType.NumField(); i++ {
		f := cl.elemType.Field(i)
		if f.PkgPath == "" { // exported
			targets[util.FldToCol(f.Name)] = i
		}
	}

	cl.colToFld = map[int]int{}
	for iC, c := range names {
		if iF, ok := targets[c]; ok {
			cl.colToFld[iC] = iF
		}
	}
	cl.slice.Set(reflect.MakeSlice(reflect.SliceOf(cl.elemType), 0, 0))
	return len(cl.colToFld) > 0, nil
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
