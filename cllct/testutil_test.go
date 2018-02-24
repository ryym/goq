package cllct_test

import (
	"reflect"

	"github.com/ryym/goq/cllct"
	"github.com/ryym/goq/goql"
)

func sel(alias, strct, field string) goql.Selection {
	return goql.Selection{TableAlias: alias, StructName: strct, FieldName: field}
}

func execCollector(
	cllcts []cllct.Collector,
	rows [][]interface{},
	selects []goql.Selection,
	colNames []string,
) error {
	if selects == nil {
		selects = make([]goql.Selection, len(colNames))
	} else {
		colNames = make([]string, len(selects))
	}

	initConf := cllct.NewInitConf(selects, colNames)
	cllcts, err := cllct.InitCollectors(cllcts, initConf)
	if err != nil {
		return err
	}

	for _, row := range rows {
		ptrs := make([]interface{}, len(selects))
		for _, cl := range cllcts {
			cl.Next(ptrs)
		}
		for i, p := range ptrs {
			if p != nil {
				reflect.ValueOf(p).Elem().Set(reflect.ValueOf(row[i]))
			}
		}
		for _, cl := range cllcts {
			cl.AfterScan(ptrs)
		}
	}

	return nil
}
