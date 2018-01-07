package gnr

import "github.com/ryym/goq/query"

func lift(v interface{}) *Ops {
	exp, ok := v.(*Ops)
	if ok {
		return exp
	}
	qr, ok := v.(query.Queryable)
	if !ok {
		qr = &litExpr{v}
	}
	return &Ops{qr}
}
