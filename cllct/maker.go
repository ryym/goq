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
	var pkFieldName string
	for _, col := range cm.cols {
		if meta := col.Meta(); meta.PK {
			pkFieldName = col.FieldName()
		}
	}
	return &ModelUniqSliceCollector{
		structName:  cm.structName,
		tableAlias:  cm.tableAlias,
		pkFieldName: pkFieldName,
		slice:       reflect.ValueOf(slice).Elem(),
		cols:        cm.cols,
	}
}
