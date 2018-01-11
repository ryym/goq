package prot

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID   int
	Name string
}

type Post struct {
	ID     int
	UserID int
}

type Pref struct {
	ID   int `goq:"pk"`
	Name string
}

type City struct {
	ID     int
	Name   string
	PrefID string
}

func Play() {
	z := NewGQL()
	fmt.Println(z)

	fmt.Println(
		z.Users.Name.Add(1).Query(),
		z.Users.ID.As("test").Query(),
		z.Users.Name.Eq("hello").Query(),
		z.Parens(z.Users.Name.Eq("hello")).As("f").Query(),
		z.Select(
			z.Parens(z.Select(z.Users.ID.Eq(2))),
			z.Users.ID,
		).Query(),
	)

	u := z.Users.As("u")
	p := z.Posts.As("p")
	query := z.Select(
		z.Users.ID,
		u.Name,
	).From(
		z.Users,
		u,
	).Where(
		u.ID.Eq(30),
		u.Name.Eq(40),
	).Joins(
		z.Users.Posts(p).Inner(),
	)

	fmt.Println(query.Query())

	fmt.Println("-------------------------")

	goDb, err := sql.Open("sqlite3", "prot/schema/prot.db")
	chk(err)
	defer goDb.Close()

	db := &DB{goDb}

	query = z.Select(z.Prefs.All()).From(z.Prefs)

	var users []User
	var prefs []Pref
	db.Query(query).Collect(
		z.Prefs.ToSlice(&prefs),
		z.Users.ToSlice(&users),
	)
	fmt.Println("RET", len(prefs), users)

	query = z.Select(z.Prefs.All(), z.Cities.All()).From(z.Prefs).Joins(
		z.InnerJoin(z.Cities).On(z.Cities.PrefID.Eq(z.Prefs.ID)),
	)

	// var cities []City
	var citiesM map[int][]City

	db.Query(query).Collect(
		z.Prefs.ToUniqSlice(&prefs),
		// z.Cities.ToSlice(&cities),
		z.Cities.ToSliceMapBy(z.Prefs.ID, &citiesM),
	)
	// fmt.Println(prefs[10], citiesM[prefs[10].ID])
	// fmt.Println(len(prefs), len(cities))
	// fmt.Println(prefs)
	// fmt.Println(cities[0:5], cities[1000:1010])

	var foos []Foo
	db.Query(
		z.Select(
			z.Prefs.ID.As("Ab"),
			z.Prefs.Name.As("Cd"),
		).From(z.Prefs).Limit(5),
	).Collect(
		z.ToSlice(&foos),
	)

	fmt.Printf("%+v\n", foos)
}

type Foo struct {
	Ab int
	Cd string
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
