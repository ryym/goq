package gnr2

import "reflect"

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

func (m *SliceCollectorMaker) ToUniqSlice(slice interface{}) *UniqSliceCollector {
	return &UniqSliceCollector{
		itemType:   m.itemType,
		structName: m.structName,
		tableAlias: m.tableAlias,
		slice:      reflect.ValueOf(slice).Elem(),
		cols:       m.cols,

		// XXX: 実際には Maker に持たせるか、itemType から抽出する。
		pkFieldName: "ID",
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
	sc.colToFld = map[int]int{}
	for iC, c := range selects {
		if c.TableAlias != "" && c.TableAlias == sc.tableAlias || c.StructName == sc.structName {
			for iF, f := range sc.cols {
				if c.FieldName == f.FieldName() {
					sc.colToFld[iC] = iF
				}
			}
		}
	}
	sc.slice.Set(reflect.MakeSlice(reflect.SliceOf(sc.itemType), 0, 0))
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

type fieldWithIdx struct {
	idx   int
	field *reflect.Value
}

type UniqSliceCollector struct {
	itemType    reflect.Type
	cols        []Column
	structName  string
	tableAlias  string
	colToFld    map[int]*fieldWithIdx
	slice       reflect.Value
	pkFieldName string
	keyIdx      int
	pks         map[interface{}]bool
	item        *reflect.Value
}

func (sc *UniqSliceCollector) Init(selects []SelectItem, _names []string) bool {
	sc.colToFld = map[int]*fieldWithIdx{}
	sc.keyIdx = -1

	// UniqSlice は毎行 struct を作る必要がないはずなので、
	// Rows.Scan()に渡すポインタ用の struct を1つだけ用意しておき、
	// 必要な行だけコピーして結果の slice に追加する。
	item := reflect.New(sc.itemType).Elem()
	sc.item = &item

	for iC, c := range selects {
		if c.TableAlias != "" && c.TableAlias == sc.tableAlias || c.StructName == sc.structName {
			if sc.pkFieldName == c.FieldName {
				sc.keyIdx = iC
			}
			for iF, f := range sc.cols {
				if c.FieldName == f.FieldName() {
					field := sc.item.Field(iF)
					sc.colToFld[iC] = &fieldWithIdx{iF, &field}
				}
			}
		}
	}

	if sc.keyIdx == -1 {
		panic("[UniqSliceCollector] primary key not found") // should return error
	}

	sc.pks = map[interface{}]bool{}
	sc.slice.Set(reflect.MakeSlice(reflect.SliceOf(sc.itemType), 0, 0))

	return len(sc.colToFld) > 0
}

func (sc *UniqSliceCollector) Next(ptrs []interface{}) {
	for c, f := range sc.colToFld {
		ptrs[c] = f.field.Addr().Interface()
	}
}

func (sc *UniqSliceCollector) AfterScan(ptrs []interface{}) {
	pk := reflect.ValueOf(ptrs[sc.keyIdx]).Elem().Interface()
	if sc.pks[pk] {
		return
	}

	copy := reflect.New(sc.itemType).Elem()
	for _, f := range sc.colToFld {
		copy.Field(f.idx).Addr().Elem().Set(*f.field)
	}

	sc.slice.Set(reflect.Append(sc.slice, copy))
	sc.pks[pk] = true
}
