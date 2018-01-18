package sample

import (
	"github.com/ryym/goq/gen/sample/models"
	"github.com/ryym/goq/gql"
)

type Users struct {
	model models.User
	name  string
	alias string

	ID   gql.Column
	Name gql.Column
}

func NewUsers() *Users {
	cm := gql.NewColumnMaker("users", "User")
	return &Users{
		model: models.User{},
		name:  "users",

		ID:   cm.Col("ID", "id"),
		Name: cm.Col("Name", "name"),
	}
}

func (t *Users) TableName() string     { return t.name }
func (t *Users) TableAlias() string    { return t.alias }
func (t *Users) All() gql.ExprListExpr { return gql.AllCols(t.Columns()) }
func (t *Users) Columns() []gql.Column {
	return []gql.Column{t.ID, t.Name}
}
func (t *Users) As(alias string) *Users {
	t2 := *t
	t2.alias = alias
	gql.CopyTableAs(alias, t, &t2)
	return &t2
}

type Cities struct {
	model City
	name  string
	alias string

	ID           gql.Column
	Name         gql.Column
	PrefectureID gql.Column
}

func NewCities() *Cities {
	cm := gql.NewColumnMaker("cities", "City")
	return &Cities{
		model: City{},
		name:  "cities",

		ID:           cm.Col("ID", "id"),
		Name:         cm.Col("Name", "name"),
		PrefectureID: cm.Col("PrefectureID", "prefecture_id"),
	}
}

func (t *Cities) TableName() string     { return t.name }
func (t *Cities) TableAlias() string    { return t.alias }
func (t *Cities) All() gql.ExprListExpr { return gql.AllCols(t.Columns()) }
func (t *Cities) Columns() []gql.Column {
	return []gql.Column{t.ID, t.Name, t.PrefectureID}
}
func (t *Cities) As(alias string) *Cities {
	t2 := *t
	t2.alias = alias
	gql.CopyTableAs(alias, t, &t2)
	return &t2
}
