package prot

import "reflect"

// XXX: 面倒だけど、テーブルは個別に定義するしかなさそう。

type JoinDef struct {
	table Table
	on    PredExpr
}

func (jd *JoinDef) Inner() JoinOn {
	return JoinOn{jd.table, jd.on, JOIN_INNER}
}

func copyTableAs(alias string, src Table, dest Table, model interface{}) {
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
		case "CollectorMaker":
			cm := NewCollectorMaker(model, dest.Columns(), alias)
			destV.Field(i).Set(reflect.ValueOf(cm))
		}
	}
}
