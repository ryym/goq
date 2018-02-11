package cllct

import (
	"fmt"
	"reflect"

	"github.com/ryym/goq/gql"
)

// For map collectors.
func isKeyCol(col *gql.Selection, key *gql.Selection) bool {
	if col.Alias != "" && col.Alias == key.Alias {
		return true
	}
	return col.StructName != "" && col.StructName == key.StructName &&
		col.TableAlias == key.TableAlias && col.FieldName == key.FieldName
}

func isSameTable(col gql.Selection, tbl tableInfo) bool {
	return col.TableAlias == tbl.tableAlias && col.StructName == tbl.structName
}

func checkPtrKind(ptr interface{}, kind reflect.Kind) error {
	tp := reflect.TypeOf(ptr)
	if tp.Kind() != reflect.Ptr || tp.Elem().Kind() != kind {
		return fmt.Errorf("required: pointer of %s, got: %s", kind, tp)
	}
	return nil
}

func checkSliceMapPtrKind(ptr interface{}) error {
	tp := reflect.TypeOf(ptr)
	if tp.Kind() == reflect.Ptr {
		elem := tp.Elem()
		if elem.Kind() == reflect.Map && elem.Elem().Kind() == reflect.Slice {
			return nil
		}
	}
	return fmt.Errorf("required: pointer of a map of slices, got: %s", tp)
}
