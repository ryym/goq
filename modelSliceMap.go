package goq

import (
	"errors"
	"reflect"
)

// ModelSliceMapCollector collects rows into a map of slices.
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
type ModelSliceMapCollector struct {
	elemType reflect.Type
	cols     []*Column
	table    tableInfo
	colToFld map[int]int
	key      Selectable
	keyIdx   int
	keyStore reflect.Value
	ptr      interface{}
	mp       reflect.Value
	row      reflect.Value
}

func (cl *ModelSliceMapCollector) ImplListCollector() {}

func (cl *ModelSliceMapCollector) init(conf *initConf) (bool, error) {
	if err := checkSliceMapPtrKind(cl.ptr); err != nil {
		return false, err
	}
	cl.mp = reflect.ValueOf(cl.ptr).Elem()
	cl.ptr = nil

	cl.colToFld = map[int]int{}
	key := cl.key.Selection()
	cl.keyIdx = -1

	for iC, c := range conf.Selects {
		if conf.canTake(iC) && isSameTable(c, cl.table) {
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

	mapType := cl.mp.Type()

	sliceType := mapType.Elem()
	cl.elemType = sliceType.Elem()
	if cl.elemType.Kind() != reflect.Struct {
		return false, errors.New("slice elem type must be struct")
	}
	cl.mp.Set(reflect.MakeMap(reflect.MapOf(mapType.Key(), sliceType)))

	return len(cl.colToFld) > 0, nil
}

func (cl *ModelSliceMapCollector) afterinit(conf *initConf) error {
	if conf.canTake(cl.keyIdx) && !cl.keyStore.IsValid() {
		return errors.New(mapKeyNotSelectedErrMsg)
	}
	return nil
}

func (cl *ModelSliceMapCollector) next(ptrs []interface{}) {
	row := reflect.New(cl.elemType).Elem()
	cl.row = row.Addr()
	for c, f := range cl.colToFld {
		ptrs[c] = row.Field(f).Addr().Interface()
	}
	if cl.keyStore.IsValid() {
		ptrs[cl.keyIdx] = cl.keyStore.Addr().Interface()
	}
}

func (cl *ModelSliceMapCollector) afterScan(ptrs []interface{}) {
	key := reflect.ValueOf(ptrs[cl.keyIdx]).Elem()
	slice := cl.mp.MapIndex(key)
	if !slice.IsValid() {
		slice = reflect.MakeSlice(reflect.SliceOf(cl.elemType), 0, 0)
	}
	slice = reflect.Append(slice, cl.row.Elem())
	cl.mp.SetMapIndex(key, slice)
}
