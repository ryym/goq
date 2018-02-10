package cllct

import "github.com/ryym/goq/gql"

func isSameTable(col gql.Selection, tbl tableInfo) bool {
	return col.TableAlias == tbl.tableAlias && col.StructName == tbl.structName
}
