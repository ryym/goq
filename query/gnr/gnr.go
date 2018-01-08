package gnr

import (
	"database/sql"
	"fmt"

	"github.com/ryym/goq/query"
)

var g *GoqGnr = &GoqGnr{}

type User struct {
	ID   int
	Name string
}

func Play() {
	q := g.Parens(g.Val(10).Add(20)).Mlt(3)
	fmt.Println(q.ToQuery())

	fmt.Println(
		g.Func("Abc", 1, 2, "test").ToQuery(),
		g.Col("users", "name").Eq(g.Col("foo", "bar").Add(1)).ToQuery(),
		g.And(
			g.Col("users", "name").Eq("NAME"),
			g.Col("users", "id").Eq(10),
		).ToQuery(),
	)

	Users := &UsersTable{
		name: "users",
		ID:   &Ops{&columnExpr{"users", "", "id", "User", "ID"}},
		Name: &Ops{&columnExpr{"users", "", "name", "User", "Name"}},
	}
	Users.SliceCollectorMaker = NewSliceCollectorMaker(User{}, Users.Columns(), "")

	u := Users.As("u")
	// fmt.Println(Users.ID)
	// fmt.Println(u.ID)

	Posts := &PostsTable{
		name:   "posts",
		ID:     &Ops{&columnExpr{"posts", "", "id", "PostsTable", "ID"}},
		UserID: &Ops{&columnExpr{"posts", "", "user_id", "PostsTable", "UserID"}},
	}
	p := Posts.As("p")

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
	).From(
		Users, u,
	).Joins(
		g.InnerJoin(p).On(Users.ID.Eq(Posts.ID)),
		Users.Posts(p).Inner(),
		u.Posts(Posts).Inner(),
	).Where(
		Users.Name.Eq("wow"),
		p.UserID.Eq(10),
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

	fmt.Println("---------------------------")

	db, err := sql.Open("sqlite3", "prot/prot.db")
	chk(err)
	defer db.Close()

	// rows, err := db.Query("select name from users")
	query := g.Select(Users.Name).From(Users)
	// fmt.Println(query.GetSelects())
	qr := query.ToQuery()
	rows, err := db.Query(qr.Query)
	chk(err)
	defer rows.Close()

	var users []User
	var collector Collector = Users.ToSlice(&users)
	collector.Init(query.GetSelects())

	rowsCols, err := rows.Columns()
	chk(err)
	ptrs := make([]interface{}, len(rowsCols))
	for rows.Next() {
		// for _, cl := range colls {
		// 	cl.Next(ptrs)
		// }

		collector.Next(ptrs)
		rows.Scan(ptrs...)
		collector.AfterScan(ptrs)

		// var name string
		// rows.Scan(&name)
		// fmt.Println(name)

		// for _, cl := range colls {
		// 	cl.AfterScan(ptrs)
		// }

		// for _, p := range ptrs {
		// 	fmt.Println(reflect.Indirect(reflect.ValueOf(p)).Interface())
		// }
	}

	fmt.Println("RET", users)
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
