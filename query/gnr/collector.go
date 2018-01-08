package gnr

import (
	"reflect"

	q "github.com/ryym/goq/query"
)

type Collector interface {
	Init(selects []q.SelectItem)
	Next(ptrs []interface{})
	AfterScan(ptrs []interface{})
}

// ひとまずモデル専用の slice collector

// テーブル struct の値を初期化するタイミングで作る
type SliceCollectorMaker struct {
	itemType   reflect.Type
	structName string
	tableAlias string
	cols       []*columnExpr
}

func NewSliceCollectorMaker(model interface{}, cols []*columnExpr, alias string) *SliceCollectorMaker {
	itemType := reflect.TypeOf(User{})
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
	cols       []*columnExpr
	structName string
	tableAlias string
	colToFl    map[int]int
	slice      reflect.Value
	item       reflect.Value
}

func (sc *SliceCollector) Init(selects []q.SelectItem) {
	colToFl := map[int]int{}
	for iC, c := range selects {
		if c.TableAlias == sc.tableAlias || c.StructName == sc.structName {
			for iF, f := range sc.cols {
				if c.FieldName == f.fieldName {
					colToFl[iC] = iF
				}
			}
		}
	}
	sc.colToFl = colToFl
}

func (sc *SliceCollector) Next(ptrs []interface{}) {
	item := reflect.New(sc.itemType).Elem()
	sc.item = item.Addr()

	for c, f := range sc.colToFl {
		ptrs[c] = item.Field(f).Addr().Interface()
	}
}

func (sc *SliceCollector) AfterScan(_ptrs []interface{}) {
	sc.slice.Set(reflect.Append(sc.slice, sc.item.Elem()))
}
