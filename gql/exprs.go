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
