package cllct_test

import "github.com/ryym/goq/gql"

func sel(alias, strct, field string) gql.Selection {
	return gql.Selection{TableAlias: alias, StructName: strct, FieldName: field}
}

//go:generate goq -testpkg gql_test.go

type Tables struct {
	users     User
	countries Country
	cities    City
	addresses Address
}

type User struct {
	ID   int `goq:"pk"`
	Name string
}

type Country struct {
	ID   int `goq:"pk"`
	Name string
}

type City struct {
	ID        int `goq:"pk"`
	Name      string
	CountryID int
}

type Address struct {
	ID         int `goq:"pk"`
	Address1   string
	Address2   string
	CityID     int
	PostalCode string
}
