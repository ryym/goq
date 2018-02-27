package cllct_test

import (
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
	return cllct.ExecCollectorsForTest(cllcts, rows, selects, colNames)
}
