package query

func Noop() {}

var g QueryGlobal

func hello(
	s SelectClause,
	p PredExpr,
	v ValExpr,
	c ColumnExpr,
	t Table,
) {
	s.Select(
		p,
		v,
		v.Add(v).As("added"),
		t.All(),
	).From(
		t,
		t.As("test"),
	).Where(
		v.Eq(v),
		g.Or(
			g.And(
				v.Gte(g.Raw(30)),
				p,
			),
			g.And(
				v.Lt(c.Mlt(c)),
				v.Between(g.Raw(10), c),
			),
		),
		p,
		p,
	).Joins(
		g.InnerJoin(t).On(c.Eq(c)),
		g.LeftJoin(t).On(c.Lt(c)),
	).GroupBy(
		c,
	).Having()
}
