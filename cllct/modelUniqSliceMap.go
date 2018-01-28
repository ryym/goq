package cllct

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/ryym/goq/gql"
)

type ModelUniqSliceMapCollector struct {
	elemType    reflect.Type
	cols        []*gql.Column
	structName  string
	tableAlias  string
	colToFld    map[int]int
	key         gql.Querier
	keyIdx      int
	keyStore    reflect.Value
	pkFieldName string
	pkIdx       int
	pks         map[interface{}]bool
	mp          reflect.Value
	elem        *reflect.Value
}

func (cl *ModelUniqSliceMapCollector) Init(selects []gql.Selection, _names []string) (bool, error) {
	cl.colToFld = map[int]int{}
	key := cl.key.Selection()
	cl.keyIdx = -1
	cl.pkIdx = -1

	for iC, c := range selects {
		if c.TableAlias == cl.tableAlias && c.StructName == cl.structName {
			if cl.pkFieldName == c.FieldName {
				cl.pkIdx = iC
			}
			for iF, f := range cl.cols {
				if c.FieldName == f.FieldName() {
					cl.colToFld[iC] = iF
				}
			}
		}

		isKey := c.Alias != "" && c.Alias == key.Alias ||
			c.TableAlias == key.TableAlias && c.StructName == key.StructName && c.FieldName == key.FieldName
		if isKey {
			cl.keyIdx = iC
		}
	}

	if cl.keyIdx == -1 {
		return false, errors.New("specified key not found")
	}
	if cl.pkIdx == -1 {
		return false, fmt.Errorf("primary key %s not selected", cl.pkFieldName)
	}

	mapType := cl.mp.Type()
	sliceType := mapType.Elem()
	cl.elemType = sliceType.Elem()
	cl.mp.Set(reflect.MakeMap(reflect.MapOf(mapType.Key(), sliceType)))

	elem := reflect.New(cl.elemType).Elem()
	cl.elem = &elem
	cl.pks = map[interface{}]bool{}

	return len(cl.colToFld) > 0, nil
}

func (cl *ModelUniqSliceMapCollector) Next(ptrs []interface{}) {
	for c, f := range cl.colToFld {
		ptrs[c] = cl.elem.Field(f).Addr().Interface()
	}
	if cl.keyStore.IsValid() {
		ptrs[cl.keyIdx] = cl.keyStore.Addr().Interface()
	}
}

func (cl *ModelUniqSliceMapCollector) AfterScan(ptrs []interface{}) {
	pk := reflect.ValueOf(ptrs[cl.pkIdx]).Elem().Interface()
	if cl.pks[pk] {
		return
	}
	cl.pks[pk] = true

	key := reflect.ValueOf(ptrs[cl.keyIdx]).Elem()
	slice := cl.mp.MapIndex(key)
	if !slice.IsValid() {
		slice = reflect.MakeSlice(reflect.SliceOf(cl.elemType), 0, 0)
	}
	slice = reflect.Append(slice, *cl.elem)
	cl.mp.SetMapIndex(key, slice)
}