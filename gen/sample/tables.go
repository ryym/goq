package sample

import "github.com/ryym/goq/gen/sample/models"

//go:generate go run "../../cmd/goq/goq.go"

type Tables struct {
	users  models.User
	cities City
}

type Common struct {
	ID int
}

type City struct {
	Common
	Name         string
	PrefectureID int
}
