package cllct

import (
	"reflect"

	"github.com/ryym/goq/gql"
)

type CollectorMaker struct{}

func NewMaker() *CollectorMaker {
	return &CollectorMaker{}
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

func (cm *ModelCollectorMaker) ToElem(elem interface{}) *ModelElemCollector {
	return &ModelElemCollector{
		structName: cm.structName,
		tableAlias: cm.tableAlias,
		elem:       reflect.ValueOf(elem).Elem(),
		cols:       cm.cols,
	}
}

func (cm *ModelCollectorMaker) ToSlice(slice interface{}) *ModelSliceCollector {
	return &ModelSliceCollector{
		structName: cm.structName,
		tableAlias: cm.tableAlias,
		slice:      reflect.ValueOf(slice).Elem(),
		cols:       cm.cols,
	}
}

func (cm *ModelCollectorMaker) ToUniqSlice(slice interface{}) *ModelUniqSliceCollector {
	pkFieldName := ""
	if pkCol := findPKCol(cm.cols); pkCol != nil {
		pkFieldName = pkCol.FieldName()
	}
	return &ModelUniqSliceCollector{
		structName:  cm.structName,
		tableAlias:  cm.tableAlias,
		pkFieldName: pkFieldName,
		slice:       reflect.ValueOf(slice).Elem(),
		cols:        cm.cols,
	}
}

func (cm *ModelCollectorMaker) ToMap(mp interface{}) *ModelMapCollector {
	return &ModelMapCollector{
		structName: cm.structName,
		tableAlias: cm.tableAlias,
		key:        findPKCol(cm.cols),
		mp:         reflect.ValueOf(mp).Elem(),
		cols:       cm.cols,
	}
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

func (cm *ModelCollectorMaker) ToSliceMap(mp interface{}) *modelSliceMapCollectorMaker {
	return &modelSliceMapCollectorMaker{&ModelSliceMapCollector{
		structName: cm.structName,
		tableAlias: cm.tableAlias,
		mp:         reflect.ValueOf(mp).Elem(),
		cols:       cm.cols,
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

func (cm *ModelCollectorMaker) ToUniqSliceMap(mp interface{}) *modelUniqSliceMapCollectorMaker {
	pkFieldName := ""
	if pkCol := findPKCol(cm.cols); pkCol != nil {
		pkFieldName = pkCol.FieldName()
	}
	return &modelUniqSliceMapCollectorMaker{&ModelUniqSliceMapCollector{
		structName:  cm.structName,
		tableAlias:  cm.tableAlias,
		pkFieldName: pkFieldName,
		mp:          reflect.ValueOf(mp).Elem(),
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
