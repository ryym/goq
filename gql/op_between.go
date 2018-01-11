package gql

import "fmt"

type betweenOp struct {
	start Querier
	end   Querier
	ops
}

func (o *betweenOp) init() *betweenOp {
	o.ops = ops{o}
	return o
}

func (o *betweenOp) Query() Query {
	s := o.start.Query()
	e := o.end.Query()
	return Query{
		fmt.Sprintf("BETWEEN %s AND %s", s.Query, e.Query),
		append(s.Args, s.Args...),
	}
}

func (o *betweenOp) Selection() Selection { return Selection{} }
