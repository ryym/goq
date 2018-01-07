package gnr

import (
	"fmt"

	"github.com/ryym/goq/query"
)

var g *GoqGnr = &GoqGnr{}

func Play() {
	q := g.Parens(g.Val(10).Add(20)).Mlt(3)
	fmt.Println(q.ToQuery())

	fmt.Println(
		g.Func("Abc", 1, 2, "test").ToQuery(),
		g.Col("users", "name").Eq(g.Col("foo", "bar").Add(1)).ToQuery(),
	)

	sl := g.Select(
		g.Col("users", "name"),
		g.Val(1),
		g.Val(1).As("n"),
		&exprListExpr{
			qs: []query.Queryable{
				g.Col("users", "id"),
				g.Col("users", "name").As("name2"),
			},
		},
		g.Col("posts", "id").Eq(1),
	)

	fmt.Println(sl.ToQuery())
	fmt.Println(sl.GetSelects())

	// 上手く行ってるけど、たったこれだけでも
	// 随分たくさんの struct をネストしないといけない。
	// predExpr{
	// 	Ops{
	// 		infixOp{
	// 			left: Ops{litExpr(10)},
	// 			right: Ops{litExpr(20)},
	// 		},
	// 	},
	// }

}
