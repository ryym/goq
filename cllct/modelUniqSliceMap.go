package cllct

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/ryym/goq/goql"
)

// ModelUniqSliceMapCollector collects rows into a map of slices.
// Each slice has uniq model structs.
// The uniqueness of a model is determined by its primary key.
//
// Example:
//
//	map[string][]City{
//		"Japan": []City{
//			{ ID: 12, Name: "Osaka" },
//			{ ID: 29, Name: "Sapporo" },
//		},
//		"Somewhere": []City{
//			{ ID: 242, Name: "Foo" },
//			{ ID: 85, Name: "Bar" },
//		},
//	}
type ModelUniqSliceMapCollector struct {
	elemType    reflect.Type
	cols        []*goql.Column
	table       tableInfo
	colToFld    map[int]int
	key         goql.Selectable
	keyIdx      int
	keyStore    reflect.Value
	pkFieldName string
	pkIdx       int
	pks         map[interface{}]bool
	ptr         interface{}
	mp          reflect.Value
	elem        *reflect.Value
}

func (cl *ModelUniqSliceMapCollector) ImplListCollector() {}

func (cl *ModelUniqSliceMapCollector) init(conf *initConf) (bool, error) {
	if err := checkSliceMapPtrKind(cl.ptr); err != nil {
		return false, err
	}
	cl.mp = reflect.ValueOf(cl.ptr).Elem()
	cl.ptr = nil

	mapType := cl.mp.Type()
	sliceType := mapType.Elem()
	cl.elemType = sliceType.Elem()
	if cl.elemType.Kind() != reflect.Struct {
		return false, errors.New("slice elem type must be struct")
	}
	cl.mp.Set(reflect.MakeMap(reflect.MapOf(mapType.Key(), sliceType)))

	cl.colToFld = map[int]int{}
	key := cl.key.Selection()
	cl.keyIdx = -1
	cl.pkIdx = -1

	for iC, c := range conf.Selects {
		if conf.canTake(iC) && isSameTable(c, cl.table) {
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

		if isKeyCol(&c, &key) {
			cl.keyIdx = iC
		}
	}

	if cl.keyIdx == -1 {
		return false, errors.New("specified key not found")
	}
	if cl.pkIdx == -1 {
		return false, fmt.Errorf("primary key %s not selected", cl.pkFieldName)
	}

	elem := reflect.New(cl.elemType).Elem()
	cl.elem = &elem
	cl.pks = map[interface{}]bool{}

	return len(cl.colToFld) > 0, nil
}

func (cl *ModelUniqSliceMapCollector) afterinit(conf *initConf) error {
	if conf.canTake(cl.keyIdx) && !cl.keyStore.IsValid() {
		return errors.New(mapKeyNotSelectedErrMsg)
	}
	return nil
}

func (cl *ModelUniqSliceMapCollector) next(ptrs []interface{}) {
	for c, f := range cl.colToFld {
		ptrs[c] = cl.elem.Field(f).Addr().Interface()
	}
	if cl.keyStore.IsValid() {
		ptrs[cl.keyIdx] = cl.keyStore.Addr().Interface()
	}
}

func (cl *ModelUniqSliceMapCollector) afterScan(ptrs []interface{}) {
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
