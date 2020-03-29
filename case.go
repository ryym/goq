package goq

// CaseExpr represents a 'CASE' expression without 'ELSE' clause.
type CaseExpr struct {
	val   Expr
	cases []*WhenClause
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

// CaseEelseExpr represents a 'CASE' expression with 'ELSE' clause.
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

// WhenClause constructs a 'WHEN' clause used for case expressions.
type WhenClause struct {
	when Expr
	then Expr
}

func (w *WhenClause) Then(then interface{}) *WhenClause {
	w.then = lift(then)
	return w
}

func (w *WhenClause) apply(q *Query, ctx DBContext) {
	q.query = append(q.query, " WHEN ")
	w.when.Apply(q, ctx)
	q.query = append(q.query, " THEN ")
	w.then.Apply(q, ctx)
}
