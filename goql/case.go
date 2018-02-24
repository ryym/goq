package goql

type CaseExpr struct {
	val   Expr
	cases []*WhenExpr
	ops
}

func (c *CaseExpr) init() *CaseExpr {
	c.ops = ops{c}
	return c
}

func (c *CaseExpr) Apply(q *Query, ctx DBContext) {
	c.apply(q, ctx)
	q.query = append(q.query, " END")
}

func (c *CaseExpr) apply(q *Query, ctx DBContext) {
	q.query = append(q.query, "CASE")
	if c.val != nil {
		q.query = append(q.query, " ")
		c.val.Apply(q, ctx)
	}
	for _, when := range c.cases {
		when.apply(q, ctx)
	}
}

func (c *CaseExpr) Selection() Selection { return Selection{} }

func (c *CaseExpr) Else(v interface{}) *CaseElseExpr {
	return (&CaseElseExpr{
		caseExpr: c,
		elseVal:  lift(v),
	}).init()
}

type CaseElseExpr struct {
	caseExpr *CaseExpr
	elseVal  Expr
	ops
}

func (ce *CaseElseExpr) init() *CaseElseExpr {
	ce.ops = ops{ce}
	return ce
}

func (ce *CaseElseExpr) Apply(q *Query, ctx DBContext) {
	ce.caseExpr.apply(q, ctx)
	q.query = append(q.query, " ELSE ")
	ce.elseVal.Apply(q, ctx)
	q.query = append(q.query, " END")
}

func (ce *CaseElseExpr) Selection() Selection { return Selection{} }

// WhenExpr does not implement Querier.
type WhenExpr struct {
	when Expr
	then Expr
}

func (w *WhenExpr) Then(then interface{}) *WhenExpr {
	w.then = lift(then)
	return w
}

func (w *WhenExpr) apply(q *Query, ctx DBContext) {
	q.query = append(q.query, " WHEN ")
	w.when.Apply(q, ctx)
	q.query = append(q.query, " THEN ")
	w.then.Apply(q, ctx)
}
