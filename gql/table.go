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
			copy := (&column{
				tableAlias: alias,
				tableName:  orig.TableName(),
				structName: orig.StructName(),
				name:       orig.ColumnName(),
				fieldName:  orig.FieldName(),
			}).init()
			destV.Field(i).Set(reflect.ValueOf(copy))
		}
	}
}

type TableHelper struct {
	name  string
	alias string
}

func (t *TableHelper) TableName() string  { return t.name }
func (t *TableHelper) TableAlias() string { return t.alias }

func (t *TableHelper) ApplyTable(q *Query, ctx DBContext) {
	q.query = append(q.query, t.name)
	if t.alias != "" {
		q.query = append(q.query, " AS ", t.alias)
	}
}

type DynmTable struct {
	name  string
	alias string
}

func (t *DynmTable) As(alias string) *DynmTable {
	t.alias = alias
	return t
}

func (t *DynmTable) ApplyTable(q *Query, ctx DBContext) {
	q.query = append(q.query, t.name)
	if t.alias != "" {
		q.query = append(q.query, " AS ", t.alias)
	}
}
