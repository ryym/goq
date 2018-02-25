package goql

type infixOp struct {
	left  Selectable
	right Selectable
	op    string
	ops
}

func (o *infixOp) init() *infixOp {
	o.ops = ops{o}
	return o
}

func (o *infixOp) Apply(q *Query, ctx DBContext) {
	o.left.Apply(q, ctx)
	q.query = append(q.query, " ", o.op, " ")
	o.right.Apply(q, ctx)
}

func (o *infixOp) Selection() Selection { return Selection{} }

type prefixOp struct {
	val Selectable
	op  string
	ops
}

func (o *prefixOp) init() *prefixOp {
	o.ops = ops{o}
	return o
}

func (o *prefixOp) Apply(q *Query, ctx DBContext) {
	q.query = append(q.query, o.op, " ")
	o.val.Apply(q, ctx)
}

func (o *prefixOp) Selection() Selection { return Selection{} }

type sufixOp struct {
	val Selectable
	op  string
	ops
}

func (o *sufixOp) init() *sufixOp {
	o.ops = ops{o}
	return o
}

func (o *sufixOp) Apply(q *Query, ctx DBContext) {
	o.val.Apply(q, ctx)
	q.query = append(q.query, " ", o.op)
}

func (o *sufixOp) Selection() Selection { return Selection{} }
