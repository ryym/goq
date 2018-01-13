package gql

import "reflect"

func AllCols(cols []Column) ExprListExpr {
	exps := make([]Expr, len(cols))
	for i, c := range cols {
		exps[i] = c
	}
	return &exprListExpr{exps}
}

func CopyTableAs(alias string, src Table, dest Table) {
	srcV := reflect.ValueOf(src).Elem()
	srcT := srcV.Type()
	destV := reflect.ValueOf(dest).Elem()
	for i := 0; i < srcT.NumField(); i++ {
		f := srcT.Field(i)
		switch f.Type.Name() {
		case "Column":
			orig := srcV.Field(i).Interface().(Column)
			copy := column{
				tableAlias: alias,
				tableName:  orig.TableName(),
				structName: orig.StructName(),
				name:       orig.ColumnName(),
				fieldName:  orig.FieldName(),
			}.init()
			destV.Field(i).Set(reflect.ValueOf(&copy))
		}
	}
}
