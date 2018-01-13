package gql

type betweenOp struct {
	start Querier
	end   Querier
	ops
}

func (o *betweenOp) init() *betweenOp {
	o.ops = ops{o}
	return o
}

func (o *betweenOp) Apply(q *Query, ctx DBContext) {
	q.query = append(q.query, "BETWEEN ")
	o.start.Apply(q, ctx)
	q.query = append(q.query, " AND ")
	o.end.Apply(q, ctx)
}

func (o *betweenOp) Selection() Selection { return Selection{} }
