package main

import (
	"fmt"

	"github.com/ryym/goq/gql"
)

func main() {
	q := gql.NewBuilder()
	fmt.Println(
		q.Var(1).Eq(2).Query(),
		q.Var(1).Gte(2).Query(),
		q.Var(1).Lt(2).Query(),
		q.Var(1).Between(0, 5).Query(),
		q.Var(3).IsNull().Query(),
	)
}
