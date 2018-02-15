package tests

//go:generate goq -test goq_test.go

type Tables struct {
	countries    Country
	cities       City
	addresses    Address
	technologies Tech `goq:"helper:Techs"`
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
	ID     int `goq:"pk"`
	Name   string
	CityID int
}

type Tech struct {
	ID   int `goq:"pk"`
	Name string
	Desc string `goq:"name:description"`
}
