package gql

import "fmt"

type logicalOp struct {
	op    string
	preds []PredExpr
	ops
}

func (l *logicalOp) ImplPredExpr() {}

func (l *logicalOp) init() *logicalOp {
	l.ops = ops{l}
	return l
}

func (l *logicalOp) Query() Query {
	if len(l.preds) == 0 {
		return Query{"", []interface{}{}}
	}

	pred := l.preds[0]
	for i := 1; i < len(l.preds); i++ {
		pred = &predExpr{(&infixOp{
			left:  pred,
			right: l.preds[i],
			op:    l.op,
		}).init()}
	}

	qr := pred.Query()
	return Query{
		fmt.Sprintf("(%s)", qr.Query),
		qr.Args,
	}
}

func (l *logicalOp) Selection() Selection { return Selection{} }
