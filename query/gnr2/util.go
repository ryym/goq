package gnr2

func lift(v interface{}) Expr {
	exp, ok := v.(Expr)
	if ok {
		return exp
	}
	return (&litExpr{val: v}).init()
}
