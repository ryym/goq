package tests

import "time"

//go:generate goq -test gql_test.go

type Tables struct {
	countries Country
	cities    City
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
