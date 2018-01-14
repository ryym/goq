package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ryym/goq/cllct"
	"github.com/ryym/goq/dialect"
	"github.com/ryym/goq/gql"
)

type UsersTable struct {
	name  string
	alias string
	ID    gql.Column
	Name  gql.Column
}

func (t *UsersTable) TableName() string     { return t.name }
func (t *UsersTable) TableAlias() string    { return t.alias }
func (t *UsersTable) All() gql.ExprListExpr { return gql.AllCols(t.Columns()) }
func (t *UsersTable) Columns() []gql.Column { return []gql.Column{t.ID, t.Name} }
func (t *UsersTable) As(alias string) *UsersTable {
	t2 := *t
	t2.alias = alias
	gql.CopyTableAs(alias, t, &t2)
	return &t2
}

type User struct {
	ID   int
	Name string
}

func main() {
	dl, err := dialect.New("postgres")
	if err != nil {
		panic(err)
	}
	q := gql.NewBuilder(dl)

	cm := gql.NewColumnMaker("users", "User")
	id := cm.Col("ID", "id")
	name := cm.Col("Name", "name")
	Users := &UsersTable{
		name:  "users",
		alias: "",
		ID:    id,
		Name:  name,
	}

	qs := []gql.Querier{
		q.Var(1).Eq(2),

		q.Var(1).Eq(2),

		q.Var(1).Gte(2),
		q.Var(1).Lt(2),
		q.Var(1).Between(0, 5),
		q.Var(3).IsNull(),

		q.Var(5).Add(id).Sbt(2),
		q.Var(8).Mlt(2).Eq(id),
		name.Concat(q.Var("hello")),
		q.Concat(name, "test", "hello"),

		q.Raw("now()").Sbt(1),
		q.Parens(q.Var(1).Add(2)).Mlt(3),

		q.And(
			id.Eq(1),
			q.Or(q.Var(1).Gte(3), q.Var(1).Lt(0)),
			q.Var(1).Eq(1),
		),
		q.Not(q.Or(
			name.Eq(id),
			q.Var(1).Eq(1),
		)),

		q.Func("foo", 1, 2).Add(3),
		q.Count(q.Var(10)),
		q.Coalesce(name, q.Var(20)),

		q.Select(id, name, q.Var(1).Add(id).As("test")),
		q.Select(id).Limit(3).OrderBy(id),
		q.Select(id, Users.All(), q.Var(1)).From(Users),
		q.Select(Users.All()).From(Users).Joins(
			q.LeftJoin(Users).On(Users.Name.Eq("bob")),
		).Where(
			Users.ID.Gte(3),
			Users.Name.Like("%bob"),
		).GroupBy(
			Users.ID,
			Users.Name,
		).Having(
			q.Count(Users.ID).Lt(100),
		).OrderBy(
			Users.ID,
		).Limit(10).Offset(20),

		q.Select(Users.ID).From(Users).Where(
			Users.ID.In(1, 2, 3),
			q.Exists(q.Select(Users.ID)),
		),

		q.Case(
			q.When(Users.ID.Eq(1)).Then(2),
			q.When(Users.ID.Eq(2)).Then(3),
		).Else(4).Add(1),
		q.CaseOf(Users.ID)(
			q.When(1).Then(2),
			q.When(2).Then(3),
		).Else(4).As("casewhen"),
	}

	for _, qr := range qs {
		fmt.Println(q.Query(qr))
	}

	fmt.Println("------------------------------------")

	db, err := Open("sqlite3", "prot/prot.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	z := db.QueryBuilder()

	rows, err := db.Query(z.Select(Users.Name).From(Users)).Rows()
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		rows.Scan(&name)
		fmt.Println(name)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	var mps []map[string]interface{}
	db.Query(z.Select(Users.Name).From(Users)).Collect(
		z.ToRowMapSlice(&mps),
	)
	fmt.Println(mps)

	var mp map[string]interface{}
	db.Query(z.Select(Users.Name).From(Users)).First(z.ToRowMap(&mp))
	fmt.Println(mp)

	var users []User
	userCllct := cllct.NewModelCollectorMaker(
		User{},
		[]gql.Column{id, name},
		"",
	)
	db.Query(
		z.Select(Users.ID, Users.Name).From(Users),
	).Collect(userCllct.ToSlice(&users))
	fmt.Println(users)
}
