package cllct

import (
	"reflect"

	"github.com/ryym/goq/gql"
)

func NewMaker() *CollectorMaker {
	return &CollectorMaker{}
}

type CollectorMaker struct{}

func (cm *CollectorMaker) ToRowMapSlice(slice *[]map[string]interface{}) *RowMapSliceCollector {
	return &RowMapSliceCollector{slice: slice}
}

func (cm *CollectorMaker) ToRowMap(mp *map[string]interface{}) *RowMapCollector {
	return &RowMapCollector{mp: mp}
}

func NewModelCollectorMaker(
	model interface{},
	cols []gql.Column,
	alias string,
) *ModelCollectorMaker {
	elemType := reflect.TypeOf(model)
	return &ModelCollectorMaker{
		elemType:   elemType,
		structName: elemType.Name(),
		tableAlias: alias,
		cols:       cols,
	}
}

type ModelCollectorMaker struct {
	elemType   reflect.Type
	structName string
	tableAlias string
	cols       []gql.Column
}

func (cm *ModelCollectorMaker) ToSlice(slice interface{}) *ModelSliceCollector {
	return &ModelSliceCollector{
		elemType:   cm.elemType,
		structName: cm.structName,
		tableAlias: cm.tableAlias,
		slice:      reflect.ValueOf(slice).Elem(),
		cols:       cm.cols,
	}
}
