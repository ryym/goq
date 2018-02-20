package gql

type UpdateMaker struct {
	table SchemaTable
	ctx   DBContext
}

func (m *UpdateMaker) Set(vals Values) *Update {
	return &Update{
		table: m.table,
		vals:  vals,
		ctx:   m.ctx,
	}
}

type Update struct {
	table SchemaTable
	vals  Values
	where Where
	ctx   DBContext
}

func (upd *Update) Where(preds ...PredExpr) *Update {
	upd.where.add(preds)
	return upd
}

func (upd *Update) Construct() Query {
	q := Query{}
	upd.Apply(&q, upd.ctx)
	return q
}

func (upd *Update) Apply(q *Query, ctx DBContext) {
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
