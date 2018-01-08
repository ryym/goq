package gnr2

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID   int
	Name string
}

var g = &Goq{}

func Play() {
	cm := ColumnMaker{"users", "", "User"}
	Users := &UsersTable{
		empModel: User{},
		name:     "users",
		ID:       cm.New("id", "ID"),
		Name:     cm.New("name", "Name"),
	}
	Users.SliceCollectorMaker = NewSliceCollectorMaker(Users.empModel, Users.Columns(), "")

	cm = ColumnMaker{"posts", "", "Post"}
	Posts := &PostsTable{
		name:   "posts",
		ID:     cm.New("id", "ID"),
		UserID: cm.New("user_id", "UserID"),
	}

	fmt.Println(
		Users.Name.Add(1),
		Users.ID.As("test").Query(),
		Users.Name.Eq("hello").Query(),
		g.Parens(Users.Name.Eq("hello")).As("f").Query(),
	)

	u := Users.As("u")
	p := Posts.As("p")
	query := g.Select(
		Users.ID,
		u.Name,
	).From(
		Users,
		u,
	).Where(
		u.ID.Eq(30),
		u.Name.Eq(40),
	).Joins(
		Users.Posts(p).Inner(),
	)

	fmt.Println(query.Query())

	fmt.Println("-------------------------")

	db, err := sql.Open("sqlite3", "prot/prot.db")
	chk(err)
	defer db.Close()

	query = g.Select(Users.Name).From(Users)
	qr := query.Query()
	rows, err := db.Query(qr.Query, qr.Args...)
	chk(err)
	defer rows.Close()

	selects := query.GetSelects()
	rowsCols, err := rows.Columns()
	chk(err)

	if len(rowsCols) != len(selects) {
		panic(fmt.Sprintf("rowsCols: %d, selects: %d", len(rowsCols), len(selects)))
	}

	var users []User
	var collector Collector = Users.ToSlice(&users)
	collector.Init(selects, rowsCols)

	ptrs := make([]interface{}, len(rowsCols))
	for rows.Next() {

		collector.Next(ptrs)
		rows.Scan(ptrs...)
		collector.AfterScan(ptrs)

	}

	fmt.Println("RET", users)
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
