package main

import (
	"fmt"

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

func main() {
	q := gql.NewBuilder()
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

		q.Select(Users.ID).From(Users).Where(Users.ID.In(1, 2, 3)),
	}

	for _, qr := range qs {
		fmt.Println(q.Query(qr))
	}
}
