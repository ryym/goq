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

	fmt.Println(
		q.Var(1).Eq(2).Query(),
		q.Var(1).Gte(2).Query(),
		q.Var(1).Lt(2).Query(),
		q.Var(1).Between(0, 5).Query(),
		q.Var(3).IsNull().Query(),

		q.Var(5).Add(id).Sbt(2).Query(),
		q.Var(8).Mlt(2).Eq(id).Query(),
		name.Concat(q.Var("hello")).Query(),

		q.Raw("now()").Sbt(1).Query(),
		q.Parens(q.Var(1).Add(2)).Mlt(3).Query(),

		q.And(
			id.Eq(1),
			q.Or(q.Var(1).Gte(3), q.Var(1).Lt(0)),
			q.Var(1).Eq(1),
		).Query(),
		q.Not(q.Or(
			name.Eq(id),
			q.Var(1).Eq(1),
		)).Query(),

		q.Func("foo", 1, 2).Add(3).Query(),
		q.Count(q.Var(10)).Query(),
		q.Coalesce(name, q.Var(20)).Query(),
	)
}
