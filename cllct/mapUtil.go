package cllct

import "github.com/ryym/goq/gql"

func isKeyCol(col *gql.Selection, key *gql.Selection) bool {
	if col.Alias != "" && col.Alias == key.Alias {
		return true
	}
	return col.StructName != "" && col.StructName == key.StructName &&
		col.TableAlias == key.TableAlias && col.FieldName == key.FieldName
}
