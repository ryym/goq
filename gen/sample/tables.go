package sample

import "github.com/ryym/goq/gen/sample/models"

//go:generate go run "../../cmd/goq/goq.go"

type Tables struct {
	users  models.User
	cities City
}

type City struct {
	ID           int
	Name         string
	PrefectureID int
}
