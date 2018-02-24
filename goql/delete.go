package goql

type Delete struct {
	table SchemaTable
	where Where
	ctx   DBContext
}

func (dlt *Delete) Construct() Query {
	q := Query{}
	dlt.Apply(&q, dlt.ctx)
	return q
}

func (dlt *Delete) Where(preds ...PredExpr) *Delete {
	dlt.where.add(preds)
	return dlt
}

func (dlt *Delete) Apply(q *Query, ctx DBContext) {
	q.query = append(q.query, "DELETE FROM ")
	dlt.table.ApplyTable(q, ctx)
	dlt.where.Apply(q, ctx)
}
