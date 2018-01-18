package models

type User struct {
	ID   int `goq:"pk"`
	Name string
}

// type City struct {
// 	ID           int
// 	Name         string
// 	PrefectureID int
// }
