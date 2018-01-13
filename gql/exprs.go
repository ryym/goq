package gql

type aliased struct {
	exp   Expr
	alias string
}

func (a *aliased) Alias() string { return a.alias }

func (a *aliased) Apply(q *Query, ctx DBContext) {
	a.exp.Apply(q, ctx)
	q.query = append(q.query, " AS ", a.alias)
}

func (a *aliased) Selection() Selection { return Selection{} }

type litExpr struct {
	val interface{}
	ops
}

func (l *litExpr) init() *litExpr {
	l.ops = ops{l}
	return l
}

// TODO: Add no placeholder version?
func (l *litExpr) Apply(q *Query, ctx DBContext) {
	q.query = append(q.query, ctx.Placeholder(q.args))
	q.args = append(q.args, l.val)
}

func (l *litExpr) Selection() Selection { return Selection{} }

type predExpr struct {
	Expr
}

func (p *predExpr) ImplPredExpr() {}

type rawExpr struct {
	sql string
	ops
}

func (r *rawExpr) init() *rawExpr {
	r.ops = ops{r}
	return r
}

func (r *rawExpr) Apply(q *Query, ctx DBContext) {
	q.query = append(q.query, r.sql)
}

func (r *rawExpr) Selection() Selection { return Selection{} }

type parensExpr struct {
	exp Expr
	ops
}

func (p *parensExpr) init() *parensExpr {
	p.ops = ops{p}
	return p
}

func (p *parensExpr) Apply(q *Query, ctx DBContext) {
	q.query = append(q.query, "(")
	p.exp.Apply(q, ctx)
	q.query = append(q.query, ")")
}

func (p *parensExpr) Selection() Selection { return p.exp.Selection() }

type inExpr struct {
	val  Expr
	exps []Expr
	ops
}

func (ie *inExpr) init() *inExpr {
	ie.ops = ops{ie}
	return ie
}

func (ie *inExpr) Apply(q *Query, ctx DBContext) {
	ie.val.Apply(q, ctx)
	q.query = append(q.query, " IN (")
	if len(ie.exps) > 0 {
		ie.exps[0].Apply(q, ctx)
		for i := 1; i < len(ie.exps); i++ {
			q.query = append(q.query, ", ")
			ie.exps[i].Apply(q, ctx)
		}
	}
	q.query = append(q.query, ")")
}

func (ie *inExpr) Selection() Selection { return Selection{} }

type funcExpr struct {
	name string
	args []Expr
	ops
}

func (f *funcExpr) init() *funcExpr {
	f.ops = ops{f}
	return f
}

func (f *funcExpr) Apply(q *Query, ctx DBContext) {
	q.query = append(q.query, f.name+"(")
	lastIdx := len(f.args) - 1
	for i, a := range f.args {
		a.Apply(q, ctx)
		if i < lastIdx {
			q.query = append(q.query, ", ")
		}
	}
	q.query = append(q.query, ")")
}

func (f *funcExpr) Selection() Selection { return Selection{} }

type exprListExpr struct {
	exps []Expr
}

func (el *exprListExpr) Apply(q *Query, ctx DBContext) {
	if len(el.exps) == 0 {
		return
	}
	el.exps[0].Apply(q, ctx)
	for i := 1; i < len(el.exps); i++ {
		q.query = append(q.query, ", ")
		el.exps[i].Apply(q, ctx)
	}
}

func (el *exprListExpr) Selection() Selection {
	panic("[INVALID] exprListExpr.Selection is called")
}

func (el *exprListExpr) Exprs() []Expr {
	return el.exps
}
