package gnr2

import (
	"reflect"
)

// ひとまずモデル専用の slice collector

type SliceCollectorMaker struct {
	itemType   reflect.Type
	structName string
	tableAlias string
	cols       []Column
}

func NewSliceCollectorMaker(model interface{}, cols []Column, alias string) *SliceCollectorMaker {
	itemType := reflect.TypeOf(model)
	return &SliceCollectorMaker{
		itemType:   itemType,
		structName: itemType.Name(),
		tableAlias: alias,
		cols:       cols,
	}
}

func (m *SliceCollectorMaker) ToSlice(slice interface{}) *SliceCollector {
	return &SliceCollector{
		itemType:   m.itemType,
		structName: m.structName,
		tableAlias: m.tableAlias,
		slice:      reflect.ValueOf(slice).Elem(),
		cols:       m.cols,
	}
}

type SliceCollector struct {
	itemType   reflect.Type
	cols       []Column
	structName string
	tableAlias string
	colToFld   map[int]int
	slice      reflect.Value
	current    reflect.Value
}

func (sc *SliceCollector) Init(selects []SelectItem, _names []string) bool {
	colToFld := map[int]int{}
	for iC, c := range selects {
		if c.TableAlias != "" && c.TableAlias == sc.tableAlias || c.StructName == sc.structName {
			for iF, f := range sc.cols {
				if c.FieldName == f.FieldName() {
					colToFld[iC] = iF
				}
			}
		}
	}
	sc.colToFld = colToFld

	return len(sc.colToFld) > 0
}

func (sc *SliceCollector) Next(ptrs []interface{}) {
	current := reflect.New(sc.itemType).Elem()
	sc.current = current.Addr()
	for c, f := range sc.colToFld {
		ptrs[c] = current.Field(f).Addr().Interface()
	}
}

func (sc *SliceCollector) AfterScan(_ptrs []interface{}) {
	sc.slice.Set(reflect.Append(sc.slice, sc.current.Elem()))
}
