package main

import (
	"fmt"

	"github.com/ryym/goq/gql"
)

func main() {
	q := gql.NewBuilder()
	fmt.Println(q.Var(10).Query())
}
