package gql

import (
	"errors"
	"reflect"
)

type Values map[*Column]interface{}

type InsertMaker struct {
	table SchemaTable
	cols  []*Column
	ctx   DBContext
}

func (m *InsertMaker) Values(vals interface{}, valsList ...interface{}) *Insert {
	colList := m.cols
	if colList == nil {
		colList = m.table.All().Columns()
	}
	cols := make(map[string]*Column, len(colList))
	for _, col := range colList {
		cols[col.FieldName()] = col
	}

	vl := make([]Values, len(valsList)+1)
	for i, vals := range append([]interface{}{vals}, valsList...) {
		valsMap, err := makeValuesMap(vals, cols)
		if err != nil {
			return &Insert{errs: []error{err}}
		}
		vl[i] = valsMap
	}

	return &Insert{
		table:    m.table,
		cols:     m.cols,
		valsList: vl,
		ctx:      m.ctx,
	}
}

func (m *InsertMaker) ValuesMap(vals Values, valsList ...Values) *Insert {
	vl := append([]Values{vals}, valsList...)
	return &Insert{
		table:    m.table,
		cols:     m.cols,
		valsList: vl,
		ctx:      m.ctx,
	}
}

type Insert struct {
	table    SchemaTable
	cols     []*Column
	valsList []Values
	errs     []error
	ctx      DBContext
}

func (ins *Insert) Construct() Query {
	q := Query{}
	ins.Apply(&q, ins.ctx)
	return q
}

func (ins *Insert) Apply(q *Query, ctx DBContext) {
	if len(ins.errs) > 0 {
		q.errs = append(q.errs, ins.errs...)
		return
	}

	q.query = append(q.query, "INSERT INTO ")
	ins.table.ApplyTable(q, ctx)

	if len(ins.cols) > 0 {
		q.query = append(q.query, " (")
		ins.cols[0].Apply(q, ctx)
		for i := 1; i < len(ins.cols); i++ {
			q.query = append(q.query, ", ")
			ins.cols[i].Apply(q, ctx)
		}
		q.query = append(q.query, ")")
	}

	q.query = append(q.query, " VALUES ")
	for vi, vals := range ins.valsList {
		q.query = append(q.query, "(")
		if len(vals) > 0 {
			cols := ins.cols
			if len(cols) == 0 {
				cols = ins.table.All().Columns()
			}

			for i, col := range cols {
				val, ok := vals[col]
				if ok {
					q.query = append(q.query, ctx.Placeholder("", q.args))
					q.args = append(q.args, val)
				} else {
					q.query = append(q.query, "NULL")
				}

				if i < len(cols)-1 {
					q.query = append(q.query, ", ")
				}
			}
		}
		q.query = append(q.query, ")")
		if vi < len(ins.valsList)-1 {
			q.query = append(q.query, ", ")
		}
	}
}

func makeValuesMap(vals interface{}, cols map[string]*Column) (Values, error) {
	tp := reflect.TypeOf(vals)
	if tp.Kind() != reflect.Struct {
		return nil, errors.New("[INSERT] Values() requires a struct")
	}

	mp := make(Values, len(cols))
	elem := reflect.ValueOf(vals)
	for i := 0; i < tp.NumField(); i++ {
		fld := tp.Field(i)
		if col, ok := cols[fld.Name]; ok {
			mp[col] = elem.Field(i).Interface()
		}
	}

	return mp, nil
}
