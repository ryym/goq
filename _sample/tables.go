package main

import "github.com/ryym/goq/_sample/models"

//go:generate go run "../cmd/goq/goq.go"

type Tables struct {
	users       models.User
	prefectures Pref `goq:"helper:Prefs"`
	cities      City
}

type Common struct {
	ID int
}

type Pref struct {
	Common
	Name string
}

type City struct {
	Common
	Name         string
	PrefectureID int
}
