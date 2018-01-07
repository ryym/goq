package query

func Noop() {}

var g QueryBuilder

func hello(
	p PredExpr,
	v Expr,
	c Expr,
	t Table,
) {
	g.Select(
		p,
		v,
		v.Add(v).As("added"),
		t.All(),
	).From(
		t,
		// t.As("test"),
	).Where(
		v.Eq(v),
		g.Or(
			g.And(
				v.Eq(30),
				p,
			),
			g.And(
				v.Eq(c.Mlt(c)),
				v.Eq(c),
			),
		),
		p,
		p,
	).Joins(
		g.InnerJoin(t).On(c.Eq(c)),
		g.LeftJoin(t).On(c.Eq(c)),
	).GroupBy(
		c,
	).Having()
}
