package cllct_test

import (
	"reflect"

	"github.com/ryym/goq/cllct"
	"github.com/ryym/goq/gql"
)

func sel(alias, strct, field string) gql.Selection {
	return gql.Selection{TableAlias: alias, StructName: strct, FieldName: field}
}

func execCollector(
	cl cllct.Collector,
	rows [][]interface{},
	selects []gql.Selection,
	colNames []string,
) {
	if selects == nil {
		selects = make([]gql.Selection, len(colNames))
	} else {
		colNames = make([]string, len(selects))
	}

	cl.Init(selects, colNames)
	for _, row := range rows {
		ptrs := make([]interface{}, len(selects))
		cl.Next(ptrs)
		for i, p := range ptrs {
			if p != nil {
				reflect.ValueOf(p).Elem().Set(reflect.ValueOf(row[i]))
			}
		}
		cl.AfterScan(ptrs)
	}
}