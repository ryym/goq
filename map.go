package goq

import (
	"errors"
	"reflect"

	"github.com/ryym/goq/util"
)

// MapCollector collects rows into a map of struct.
//
// Example:
//
//	map[int]JapanName{
//		6: JapanName{ Name: "Yamada", Kanji: "山田" },
//		43: JapanName{ Name: "Kato", Kanji: "加藤" },
//	}
type MapCollector struct {
	elemType reflect.Type
	colToFld map[int]int
	key      Selectable
	keyIdx   int
	keyStore reflect.Value
	ptr      interface{}
	mp       reflect.Value
	row      reflect.Value
}

func (cl *MapCollector) ImplListCollector() {}

func (cl *MapCollector) init(conf *initConf) (bool, error) {
	if err := checkPtrKind(cl.ptr, reflect.Map); err != nil {
		return false, err
	}

	cl.mp = reflect.ValueOf(cl.ptr).Elem()
	cl.ptr = nil

	cl.elemType = cl.mp.Type().Elem()

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
		return false, errors.New("key not found")
	}

	for iC, name := range conf.ColNames {
		if iF, ok := targets[name]; ok && conf.take(iC) {
			cl.colToFld[iC] = iF
		}
	}

	mapType := cl.mp.Type()
	cl.mp.Set(reflect.MakeMap(reflect.MapOf(mapType.Key(), cl.elemType)))

	return len(cl.colToFld) > 0, nil
}

func (cl *MapCollector) afterinit(conf *initConf) error {
	if conf.canTake(cl.keyIdx) && !cl.keyStore.IsValid() {
		return errors.New(mapKeyNotSelectedErrMsg)
	}
	return nil
}

func (cl *MapCollector) next(ptrs []interface{}) {
	row := reflect.New(cl.elemType).Elem()
	cl.row = row.Addr()
	for c, f := range cl.colToFld {
		ptrs[c] = row.Field(f).Addr().Interface()
	}
	if cl.keyStore.IsValid() {
		ptrs[cl.keyIdx] = cl.keyStore.Addr().Interface()
	}
}

func (cl *MapCollector) afterScan(ptrs []interface{}) {
	key := reflect.ValueOf(ptrs[cl.keyIdx]).Elem()
	cl.mp.SetMapIndex(key, cl.row.Elem())
}
