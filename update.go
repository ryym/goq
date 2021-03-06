package goq

import (
	"reflect"

	"github.com/pkg/errors"
)

// UpdateMaker constructs Update struct.
type UpdateMaker struct {
	table SchemaTable
	ctx   DBContext
}

// Set accepts values map.
func (m *UpdateMaker) Set(vals Values) *Update {
	return &Update{
		table: m.table,
		vals:  vals,
		ctx:   m.ctx,
	}
}

// Elem accepts an model.
// If the 'cols' are specified,
// only the fields corresponding to these columns are updated.
func (m *UpdateMaker) Elem(elem interface{}, cols ...*Column) *Update {
	var pkCol *Column
	for _, col := range m.table.All().Columns() {
		if col.Meta().PK {
			pkCol = col
			break
		}
	}
	if pkCol == nil {
		return &Update{err: errors.New("[Update] PK column required")}
	}

	vals, err := makeValuesMap(elem, makeColsMap(cols, m.table))
	if err != nil {
		return &Update{err: errors.Wrap(err, "[Update] Elm()")}
	}

	pkVal := vals[pkCol]
	if pkVal == nil {
		elemRef := reflect.ValueOf(elem)
		if elemRef.Type().Kind() == reflect.Ptr {
			elemRef = elemRef.Elem()
		}
		pkVal = elemRef.FieldByName(pkCol.FieldName()).Interface()
		if pkVal == nil {
			return &Update{err: errors.New("[Update] PK must have a value")}
		}
	}

	upd := &Update{
		table: m.table,
		vals:  vals,
		ctx:   m.ctx,
	}
	upd.Where(pkCol.Eq(pkVal))
	return upd
}

// Update constructs an 'UPDATE' statement.
type Update struct {
	table SchemaTable
	vals  Values
	where Where
	err   error
	ctx   DBContext
}

// Where appends conditions of the update target rows.
func (upd *Update) Where(preds ...PredExpr) *Update {
	upd.where.add(preds)
	return upd
}

func (upd *Update) Construct() (Query, error) {
	q := Query{}
	upd.Apply(&q, upd.ctx)
	return q, q.Err()
}

func (upd *Update) Apply(q *Query, ctx DBContext) {
	if upd.err != nil {
		q.errs = append(q.errs, upd.err)
		return
	}

	q.query = append(q.query, "UPDATE ")
	upd.table.ApplyTable(q, ctx)

	q.query = append(q.query, " SET ")

	// Iterate columns slice instead of vals map to ensure
	// listed columns are always in the same order.
	i := 0
	for _, col := range upd.table.All().Columns() {
		val, ok := upd.vals[col]
		if ok {
			q.query = append(q.query,
				ctx.QuoteIdent(col.ColumnName()),
				" = ",
				ctx.Placeholder("", q.args),
			)
			q.args = append(q.args, val)
			if i < len(upd.vals)-1 {
				q.query = append(q.query, ", ")
			}
			i++
		}
	}

	upd.where.Apply(q, ctx)
}
