package goq_test

import (
	"github.com/ryym/goq"
)

func sel(alias, strct, field string) goq.Selection {
	return goq.Selection{TableAlias: alias, StructName: strct, FieldName: field}
}

func execCollector(
	cllcts []goq.Collector,
	rows [][]interface{},
	selects []goq.Selection,
	colNames []string,
) error {
	return goq.ExecCollectorsForTest(cllcts, rows, selects, colNames)
}
