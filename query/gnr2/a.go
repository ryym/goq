package gnr2

import "fmt"

func Play() {
	g := &Goq{}

	cm := ColumnMaker{"users", "", "User"}
	Users := &UsersTable{
		name: "users",
		ID:   cm.New("id", "ID"),
		Name: cm.New("name", "Name"),
	}

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
}
