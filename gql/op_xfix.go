package gql

import "fmt"

type infixOp struct {
	left  Querier
	right Querier
	op    string
	ops
}

func (o *infixOp) init() *infixOp {
	o.ops = ops{o}
	return o
}

func (o *infixOp) Query() Query {
	l := o.left.Query()
	r := o.right.Query()
	return Query{
		fmt.Sprintf("%s %s %s", l.Query, o.op, r.Query),
		append(l.Args, r.Args...),
	}
}

func (o *infixOp) Selection() Selection { return Selection{} }

type prefixOp struct {
	val Querier
	op  string
	ops
}

func (o *prefixOp) init() *prefixOp {
	o.ops = ops{o}
	return o
}

func (o *prefixOp) Query() Query {
	qr := o.val.Query()
	return Query{
		fmt.Sprintf("%s%s", o.op, qr.Query),
		qr.Args,
	}
}

func (o *prefixOp) Selection() Selection { return Selection{} }

type sufixOp struct {
	val Querier
	op  string
	ops
}

func (o *sufixOp) init() *sufixOp {
	o.ops = ops{o}
	return o
}

func (o *sufixOp) Query() Query {
	qr := o.val.Query()
	return Query{
		fmt.Sprintf("%s %s", qr.Query, o.op),
		qr.Args,
	}
}

func (o *sufixOp) Selection() Selection { return Selection{} }
