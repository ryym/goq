package cllct

import (
	"reflect"

	"github.com/ryym/goq/gql"
)

type CollectorMaker struct{}

func NewMaker() *CollectorMaker {
	return &CollectorMaker{}
}

func (cm *CollectorMaker) ToElem(ptr interface{}) *ElemCollector {
	return &ElemCollector{
		ptr: ptr,
	}
}

func (cm *CollectorMaker) ToSlice(ptr interface{}) *SliceCollector {
	return &SliceCollector{
		ptr: ptr,
	}
}

type mapCollectorMaker struct {
	collector *MapCollector
}

func (m *mapCollectorMaker) By(key gql.Querier) *MapCollector {
	m.collector.key = key
	return m.collector
}

func (m *mapCollectorMaker) ByWith(ptr interface{}, key gql.Querier) *MapCollector {
	m.collector.key = key
	m.collector.keyStore = reflect.ValueOf(ptr).Elem()
	return m.collector
}

func (cm *CollectorMaker) ToMap(ptr interface{}) *mapCollectorMaker {
	return &mapCollectorMaker{&MapCollector{
		ptr: ptr,
	}}
}

type sliceMapCollector struct {
	collector *SliceMapCollector
}

func (m *sliceMapCollector) By(key gql.Querier) *SliceMapCollector {
	m.collector.key = key
	return m.collector
}

func (m *sliceMapCollector) ByWith(ptr interface{}, key gql.Querier) *SliceMapCollector {
	m.collector.key = key
	m.collector.keyStore = reflect.ValueOf(ptr).Elem()
	return m.collector
}

func (cm *CollectorMaker) ToSliceMap(ptr interface{}) *sliceMapCollector {
	return &sliceMapCollector{&SliceMapCollector{
		ptr: ptr,
	}}
}

func (cm *CollectorMaker) ToRowMapSlice(slice *[]map[string]interface{}) *RowMapSliceCollector {
	return &RowMapSliceCollector{slice: slice}
}

func (cm *CollectorMaker) ToRowMap(mp *map[string]interface{}) *RowMapCollector {
	return &RowMapCollector{mp: mp}
}

type ModelCollectorMaker struct {
	structName string
	tableAlias string
	cols       []*gql.Column
}

func NewModelCollectorMaker(cols []*gql.Column, alias string) *ModelCollectorMaker {
	var structName string
	if len(cols) > 0 {
		structName = cols[0].StructName()
	}
	return &ModelCollectorMaker{
		structName: structName,
		tableAlias: alias,
		cols:       cols,
	}
}

func (cm *ModelCollectorMaker) ToElem(ptr interface{}) *ModelElemCollector {
	return &ModelElemCollector{
		table: tableInfo{cm.structName, cm.tableAlias},
		ptr:   ptr,
		cols:  cm.cols,
	}
}

func (cm *ModelCollectorMaker) ToSlice(ptr interface{}) *ModelSliceCollector {
	return &ModelSliceCollector{
		table: tableInfo{cm.structName, cm.tableAlias},
		ptr:   ptr,
		cols:  cm.cols,
	}
}

func (cm *ModelCollectorMaker) ToUniqSlice(ptr interface{}) *ModelUniqSliceCollector {
	pkFieldName := ""
	if pkCol := findPKCol(cm.cols); pkCol != nil {
		pkFieldName = pkCol.FieldName()
	}
	return &ModelUniqSliceCollector{
		table:       tableInfo{cm.structName, cm.tableAlias},
		pkFieldName: pkFieldName,
		ptr:         ptr,
		cols:        cm.cols,
	}
}

func (cm *ModelCollectorMaker) ToMap(ptr interface{}) *ModelMapCollector {
	mapCllct := &ModelMapCollector{
		table: tableInfo{cm.structName, cm.tableAlias},
		ptr:   ptr,
		cols:  cm.cols,
	}
	if pkCol := findPKCol(cm.cols); pkCol != nil {
		keySel := pkCol.Selection()
		mapCllct.keySel = &keySel
	}
	return mapCllct
}

type modelSliceMapCollectorMaker struct {
	collector *ModelSliceMapCollector
}

func (m *modelSliceMapCollectorMaker) By(key gql.Querier) *ModelSliceMapCollector {
	m.collector.key = key
	return m.collector
}

// Sometimes you need to provide a pointer to store each key value of the result map.
// For example, the pattern A is OK because `Countries.ID` will be mapped to each Country model.
// But the pattern B fails because there is no place each `Country.ID` is mapped to.
// In this case, you need to use `ByWith` to provide a place for `Country.ID` mapping (pattern C).
//
// pattern A (GOOD):
//     Collect(
//         Countries.ToSlice(&countries),
//         Cities.ToSliceMap(&cities).By(Countries.ID),
//     )
//
// pattern B (BAD):
//     Collect(Cities.ToSliceMap(&cities).By(Countries.ID))
//
// pattern C (GOOD):
//     Collect(Cities.ToSliceMap(&cities).ByWith(&countryID, Countries.ID))
func (m *modelSliceMapCollectorMaker) ByWith(ptr interface{}, key gql.Querier) *ModelSliceMapCollector {
	m.collector.key = key
	m.collector.keyStore = reflect.ValueOf(ptr).Elem()
	return m.collector
}

func (cm *ModelCollectorMaker) ToSliceMap(ptr interface{}) *modelSliceMapCollectorMaker {
	return &modelSliceMapCollectorMaker{&ModelSliceMapCollector{
		table: tableInfo{cm.structName, cm.tableAlias},
		ptr:   ptr,
		cols:  cm.cols,
	}}
}

type modelUniqSliceMapCollectorMaker struct {
	collector *ModelUniqSliceMapCollector
}

func (m *modelUniqSliceMapCollectorMaker) By(key gql.Querier) *ModelUniqSliceMapCollector {
	m.collector.key = key
	return m.collector
}

// See modelSliceMapCollectorMaker.ByWith
func (m *modelUniqSliceMapCollectorMaker) ByWith(ptr interface{}, key gql.Querier) *ModelUniqSliceMapCollector {
	m.collector.key = key
	m.collector.keyStore = reflect.ValueOf(ptr).Elem()
	return m.collector
}

func (cm *ModelCollectorMaker) ToUniqSliceMap(ptr interface{}) *modelUniqSliceMapCollectorMaker {
	pkFieldName := ""
	if pkCol := findPKCol(cm.cols); pkCol != nil {
		pkFieldName = pkCol.FieldName()
	}
	return &modelUniqSliceMapCollectorMaker{&ModelUniqSliceMapCollector{
		table:       tableInfo{cm.structName, cm.tableAlias},
		pkFieldName: pkFieldName,
		ptr:         ptr,
		cols:        cm.cols,
	}}
}

func findPKCol(cols []*gql.Column) *gql.Column {
	for _, col := range cols {
		if meta := col.Meta(); meta.PK {
			return col
		}
	}
	return nil
}
