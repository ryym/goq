package cllct

import (
	"reflect"

	"github.com/ryym/goq/gql"
)

type ModelElemCollector struct {
	cols       []*gql.Column
	structName string
	tableAlias string
	colToFld   map[int]int
	elem       reflect.Value
}

func (cl *ModelElemCollector) ImplSingleCollector() {}

func (cl *ModelElemCollector) Init(conf *InitConf) (bool, error) {
	cl.colToFld = map[int]int{}
	for iC, c := range conf.Selects {
		if conf.canTake(iC) && c.TableAlias == cl.tableAlias && c.StructName == cl.structName {
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

func (cl *ModelElemCollector) Next(ptrs []interface{}) {
	for c, f := range cl.colToFld {
		ptrs[c] = cl.elem.Field(f).Addr().Interface()
	}
}

func (cl *ModelElemCollector) AfterScan(_ptrs []interface{}) {}
