package models

type User struct {
	ID   int `goq:"pk"`
	Name string
}
