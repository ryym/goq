package cllct

import (
	"errors"
	"reflect"

	"github.com/ryym/goq/gql"
)

type ModelMapCollector struct {
	elemType   reflect.Type
	cols       []*gql.Column
	structName string
	tableAlias string
	colToFld   map[int]int
	key        gql.Querier
	keyIdx     int
	mp         reflect.Value
	row        reflect.Value
}

func (cl *ModelMapCollector) ImplListCollector() {}

func (cl *ModelMapCollector) Init(selects []gql.Selection, _names []string) (bool, error) {
	if cl.key == nil {
		return false, errors.New("PK column required")
	}

	cl.colToFld = map[int]int{}
	key := cl.key.Selection()
	cl.keyIdx = -1

	for iC, c := range selects {
		if c.TableAlias == cl.tableAlias && c.StructName == cl.structName {
			for iF, f := range cl.cols {
				if c.FieldName == f.FieldName() {
					cl.colToFld[iC] = iF
				}
			}
		}

		if isKeyCol(&c, &key) {
			cl.keyIdx = iC
		}
	}

	if cl.keyIdx == -1 {
		return false, errors.New("key not found")
	}

	mapType := cl.mp.Type()

	cl.elemType = mapType.Elem()
	cl.mp.Set(reflect.MakeMap(reflect.MapOf(mapType.Key(), cl.elemType)))

	return len(cl.colToFld) > 0, nil
}

func (cl *ModelMapCollector) Next(ptrs []interface{}) {
	row := reflect.New(cl.elemType).Elem()
	cl.row = row.Addr()
	for c, f := range cl.colToFld {
		ptrs[c] = row.Field(f).Addr().Interface()
	}
}

func (cl *ModelMapCollector) AfterScan(ptrs []interface{}) {
	key := reflect.ValueOf(ptrs[cl.keyIdx]).Elem()
	cl.mp.SetMapIndex(key, cl.row.Elem())
}
