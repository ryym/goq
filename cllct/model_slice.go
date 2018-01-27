package cllct

import (
	"reflect"

	"github.com/ryym/goq/gql"
)

type ModelSliceCollector struct {
	elemType   reflect.Type
	cols       []*gql.Column
	structName string
	tableAlias string
	colToFld   map[int]int
	slice      reflect.Value
	row        reflect.Value
}

func (cl *ModelSliceCollector) ImplListCollector() {}

func (cl *ModelSliceCollector) Init(selects []gql.Selection, _names []string) (bool, error) {
	cl.colToFld = map[int]int{}
	for iC, c := range selects {
		if c.TableAlias == cl.tableAlias && c.StructName == cl.structName {
			for iF, f := range cl.cols {
				if c.FieldName == f.FieldName() {
					cl.colToFld[iC] = iF
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
