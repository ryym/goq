package main

import (
	"fmt"

	"github.com/ryym/goq/gql"
)

func main() {
	q := gql.NewBuilder()
	cm := gql.NewColumnMaker("users", "User")
	id := cm.Col("ID", "id")
	name := cm.Col("Name", "name")

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
	}

	for _, qr := range qs {
		fmt.Println(q.Query(qr))
	}
}
