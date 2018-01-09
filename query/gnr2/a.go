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

type Pref struct {
	ID   int
	Name string
}

type City struct {
	ID     int
	Name   string
	PrefID string
}

var g = &Goq{}

func Play() {
	cm := ColumnMaker{"users", "", "User"}
	Users := &UsersTable{
		empModel: User{},
		name:     "users",
		ID:       cm.Col("ID", "id"),
		Name:     cm.Col("Name", "name"),
	}
	Users.CollectorMaker = NewCollectorMaker(Users.empModel, Users.Columns(), "")

	cm = ColumnMaker{"posts", "", "Post"}
	Posts := &PostsTable{
		name:   "posts",
		ID:     cm.Col("ID", "id"),
		UserID: cm.Col("UserID", "user_id"),
	}

	cm = ColumnMaker{"prefectures", "", "Pref"}
	Prefs := &PrefsTable{
		empModel: Pref{},
		name:     "prefectures",
		ID:       cm.Col("ID", "id"),
		Name:     cm.Col("Name", "name"),
	}
	Prefs.CollectorMaker = NewCollectorMaker(Prefs.empModel, Prefs.Columns(), "")

	cm = ColumnMaker{"cities", "", "City"}
	Cities := &CitiesTable{
		empModel: City{},
		name:     "cities",
		ID:       cm.Col("ID", "id"),
		Name:     cm.Col("Name", "name"),
		PrefID:   cm.Col("PrefID", "prefecture_id"),
	}
	Cities.CollectorMaker = NewCollectorMaker(Cities.empModel, Cities.Columns(), "")

	dest := *Users
	copyTableAs("uu", Users, &dest)

	fmt.Println(
		Users.Name.Add(1).Query(),
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

	goDb, err := sql.Open("sqlite3", "prot/prot.db")
	chk(err)
	defer goDb.Close()

	db := &DB{goDb}

	query = g.Select(Prefs.All()).From(Prefs)

	var users []User
	var prefs []Pref
	db.Query(query).Collect(
		Prefs.ToSlice(&prefs),
		Users.ToSlice(&users),
	)
	fmt.Println("RET", len(prefs), users)

	query = g.Select(Prefs.All(), Cities.All()).From(Prefs).Joins(
		g.InnerJoin(Cities).On(Cities.PrefID.Eq(Prefs.ID)),
	)

	// var cities []City
	var citiesM map[int][]City

	db.Query(query).Collect(
		Prefs.ToUniqSlice(&prefs),
		// Cities.ToSlice(&cities),
		Cities.ToSliceMapBy(Prefs.ID, &citiesM),
	)
	fmt.Println(prefs[10], citiesM[prefs[10].ID])
	// fmt.Println(len(prefs), len(cities))
	// fmt.Println(prefs)
	// fmt.Println(cities[0:5], cities[1000:1010])
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
