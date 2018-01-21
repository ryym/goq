package main

import (
	"fmt"

	"github.com/ryym/goq/dialect"
)

func main() {
	z := NewBuilder(dialect.New("postgres"))
	u := z.Users.As("u")
	fmt.Println(z.Query(z.Select(z.Users.All(), u.All()).From(z.Users, u)))
}
