package goq

import (
	"reflect"

	"github.com/pkg/errors"
)

type Values map[*Column]interface{}

// InsertMaker constructs Insert struct.
type InsertMaker struct {
	table SchemaTable
	cols  []*Column
	ctx   DBContext
}

// Values accepts one or more model structs to be inserted.
func (m *InsertMaker) Values(elems ...interface{}) *Insert {
	if len(elems) == 0 {
		panic("[goq] elems slice is empty")
	}

	cols := makeColsMap(m.cols, m.table)

	valsList := make([]Values, 0, len(elems))
	for _, elem := range elems {
		vals, err := makeValuesMap(elem, cols)
		if err != nil {
			return &Insert{err: errors.Wrap(err, "[Insert] Values()")}
		}
		valsList = append(valsList, vals)
	}

	return &Insert{
		table:    m.table,
		cols:     m.cols,
		valsList: valsList,
		ctx:      m.ctx,
	}
}

// ValuesMap accepts one or more value maps to be inserted.
// If the target columns are specified by QueryBuilder.InsertInto,
// values of non-target columns are ignored.
func (m *InsertMaker) ValuesMap(valsList ...Values) *Insert {
	if len(valsList) == 0 {
		panic("[goq] values slice is empty")
	}
	return &Insert{
		table:    m.table,
		cols:     m.cols,
		valsList: valsList,
		ctx:      m.ctx,
	}
}

// Insert constructs an 'INSERT' statement.
type Insert struct {
	table    SchemaTable
	cols     []*Column
	valsList []Values
	err      error
	ctx      DBContext
}

func (ins *Insert) Construct() (Query, error) {
	q := Query{}
	ins.Apply(&q, ins.ctx)
	return q, q.Err()
}

func (ins *Insert) Apply(q *Query, ctx DBContext) {
	if ins.err != nil {
		q.errs = append(q.errs, ins.err)
		return
	}

	q.query = append(q.query, "INSERT INTO ")
	ins.table.ApplyTable(q, ctx)

	if len(ins.cols) > 0 {
		q.query = append(q.query, " (")
		for i, col := range ins.cols {
			q.query = append(q.query, ctx.QuoteIdent(col.ColumnName()))
			if i < len(ins.cols)-1 {
				q.query = append(q.query, ", ")
			}
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

func makeColsMap(cols []*Column, table SchemaTable) map[string]*Column {
	colList := cols
	if colList == nil {
		colList = table.All().Columns()
	}
	mp := make(map[string]*Column, len(colList))
	for _, col := range colList {
		mp[col.FieldName()] = col
	}
	return mp
}

func makeValuesMap(elem interface{}, cols map[string]*Column) (Values, error) {
	tp := reflect.TypeOf(elem)
	var elemRfl reflect.Value

	if tp.Kind() == reflect.Ptr {
		elemRfl = reflect.ValueOf(elem).Elem()
		tp = elemRfl.Type()
	} else if tp.Kind() == reflect.Struct {
		elemRfl = reflect.ValueOf(elem)
	}

	if !elemRfl.IsValid() {
		return nil, errors.New("elem is not a struct nor a pointer to struct")
	}

	mp := make(Values, len(cols))
	for i := 0; i < tp.NumField(); i++ {
		fld := tp.Field(i)
		if col, ok := cols[fld.Name]; ok {
			mp[col] = elemRfl.Field(i).Interface()
		}
	}

	return mp, nil
}
