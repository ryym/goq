package cllct

import (
	"reflect"

	"github.com/ryym/goq/gql"
	"github.com/ryym/goq/util"
)

type ElemCollector struct {
	colToFld map[int]int
	elem     reflect.Value
}

func (cl *ElemCollector) ImplSingleCollector() {}

func (cl *ElemCollector) Init(selects []gql.Selection, names []string) (bool, error) {
	targets := map[string]int{}
	elemType := cl.elem.Type()
	for i := 0; i < elemType.NumField(); i++ {
		f := elemType.Field(i)
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
	return len(cl.colToFld) > 0, nil
}

func (cl *ElemCollector) Next(ptrs []interface{}) {
	for c, f := range cl.colToFld {
		ptrs[c] = cl.elem.Field(f).Addr().Interface()
	}
}

func (cl *ElemCollector) AfterScan(_ptrs []interface{}) {}
