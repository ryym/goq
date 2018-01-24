package main

import "github.com/ryym/goq/_sample/models"

//go:generate goq

type Tables struct {
	users       models.User
	prefectures Pref `goq:"helper:Prefs"`
	cities      City
}

type Common struct {
	ID int `goq:"pk"`
}

type Pref struct {
	Common
	Name string
}

type City struct {
	Common
	Name   string
	PrefID int  `goq:"name:prefecture_id"`
	Foo    bool `goq:"-"`
}
