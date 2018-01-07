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

	Users := &UsersTable{
		name: "users",
		ID:   &Ops{&columnExpr{"users", "", "id", "UsersTable", "ID"}},
		Name: &Ops{&columnExpr{"users", "", "name", "UsersTable", "Name"}},
	}
	u := Users.As("u")
	// fmt.Println(Users.ID)
	// fmt.Println(u.ID)

	sl := g.Select(
		g.Col("users", "name"),
		g.Val(1),
		g.Val(1).As("n"),
		Users.ID,
		u.ID,
		u.ID.Add(Users.ID),
		&exprListExpr{
			qs: []query.Queryable{
				Users.ID,
				Users.Name,
			},
		},
		g.Col("posts", "id").Eq(1),
	).From(Users, u)

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
