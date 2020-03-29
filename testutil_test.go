package goq_test

import (
	"github.com/ryym/goq"
	"github.com/ryym/goq/goql"
)

func sel(alias, strct, field string) goql.Selection {
	return goql.Selection{TableAlias: alias, StructName: strct, FieldName: field}
}

func execCollector(
	cllcts []goq.Collector,
	rows [][]interface{},
	selects []goql.Selection,
	colNames []string,
) error {
	return goq.ExecCollectorsForTest(cllcts, rows, selects, colNames)
}
