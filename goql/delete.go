package goql

// Delete constructs a 'DELETE' statement.
type Delete struct {
	table SchemaTable
	where Where
	ctx   DBContext
}

// Where appends conditions of the deletion target rows.
func (dlt *Delete) Where(preds ...PredExpr) *Delete {
	dlt.where.add(preds)
	return dlt
}

func (dlt *Delete) Construct() (Query, error) {
	q := Query{}
	dlt.Apply(&q, dlt.ctx)
	return q, q.Err()
}

func (dlt *Delete) Apply(q *Query, ctx DBContext) {
	q.query = append(q.query, "DELETE FROM ")
	dlt.table.ApplyTable(q, ctx)
	dlt.where.Apply(q, ctx)
}
