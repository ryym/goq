package cllct

import (
	"errors"
	"reflect"

	"github.com/ryym/goq/goql"
	"github.com/ryym/goq/util"
)

// SliceMapCollector collects rows into a map of slices.
//
// Example:
//
//	map[string][]IDs{
//		"Japan": []IDs{
//			{ CountryID: 1, CityID: 3 },
//			{ CountryID: 1, CityID: 8 },
//		},
//		"Somewhere": []IDs{
//			{ CountryID: 2, CityID: 55 },
//		},
//	}
type SliceMapCollector struct {
	elemType reflect.Type
	colToFld map[int]int
	key      goql.Selectable
	keyIdx   int
	keyStore reflect.Value
	ptr      interface{}
	mp       reflect.Value
	row      reflect.Value
}

func (cl *SliceMapCollector) ImplListCollector() {}

func (cl *SliceMapCollector) init(conf *initConf) (bool, error) {
	if err := checkSliceMapPtrKind(cl.ptr); err != nil {
		return false, err
	}

	cl.mp = reflect.ValueOf(cl.ptr).Elem()
	cl.elemType = cl.mp.Type().Elem().Elem()
	cl.ptr = nil

	mapType := cl.mp.Type()
	cl.mp.Set(reflect.MakeMap(reflect.MapOf(mapType.Key(), mapType.Elem())))

	targets := map[string]int{}
	for i := 0; i < cl.elemType.NumField(); i++ {
		f := cl.elemType.Field(i)
		if f.PkgPath == "" { // exported
			targets[util.FldToCol(f.Name)] = i
		}
	}

	cl.colToFld = map[int]int{}
	key := cl.key.Selection()
	cl.keyIdx = -1

	for iC, name := range conf.ColNames {
		if name == key.Alias || isKeyCol(&conf.Selects[iC], &key) {
			cl.keyIdx = iC
		}
	}
	if cl.keyIdx == -1 {
		return false, errors.New("specified key not found")
	}

	for iC, name := range conf.ColNames {
		if iF, ok := targets[name]; ok && conf.take(iC) {
			cl.colToFld[iC] = iF
		}
	}

	return len(cl.colToFld) > 0, nil
}

func (cl *SliceMapCollector) afterinit(conf *initConf) error {
	if conf.canTake(cl.keyIdx) && !cl.keyStore.IsValid() {
		return errors.New(mapKeyNotSelectedErrMsg)
	}
	return nil
}

func (cl *SliceMapCollector) next(ptrs []interface{}) {
	row := reflect.New(cl.elemType).Elem()
	cl.row = row.Addr()
	for c, f := range cl.colToFld {
		ptrs[c] = row.Field(f).Addr().Interface()
	}
	if cl.keyStore.IsValid() {
		ptrs[cl.keyIdx] = cl.keyStore.Addr().Interface()
	}
}

func (cl *SliceMapCollector) afterScan(ptrs []interface{}) {
	key := reflect.ValueOf(ptrs[cl.keyIdx]).Elem()
	slice := cl.mp.MapIndex(key)
	if !slice.IsValid() {
		slice = reflect.MakeSlice(reflect.SliceOf(cl.elemType), 0, 0)
	}
	slice = reflect.Append(slice, cl.row.Elem())
	cl.mp.SetMapIndex(key, slice)
}
