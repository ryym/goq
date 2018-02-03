package tests

import "time"

//go:generate goq -test gql_test.go

type Tables struct {
	countries    Country
	cities       City
	addresses    Address
	technologies Tech `goq:"helper:Techs"`
}

type Country struct {
	ID        int `goq:"pk"`
	Name      string
	UpdatedAt time.Time
}

type City struct {
	ID        int `goq:"pk"`
	Name      string
	CountryID int
	UpdatedAt time.Time
}

type Address struct {
	ID        int `goq:"pk"`
	Name      string
	CityID    int
	UpdatedAt time.Time
}

type Tech struct {
	ID   int `goq:"pk"`
	Name string
	Desc string `goq:"name:description"`
}
