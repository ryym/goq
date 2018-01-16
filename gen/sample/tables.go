package sample

//go:generate go run "../../main.go"

type Tables struct {
	users  User
	cities City
}

type User struct {
	ID   int `goq:"pk"`
	Name string
}

type City struct {
	ID           int
	Name         string
	PrefectureID int
}
