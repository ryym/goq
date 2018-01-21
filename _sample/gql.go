package main

import (
	"github.com/ryym/goq"
	"github.com/ryym/goq/cllct"
	"github.com/ryym/goq/dialect"
	"github.com/ryym/goq/gql"
)

type Users struct {
	gql.Table
	*cllct.ModelCollectorMaker

	ID   gql.Column
	Name gql.Column
}

func NewUsers(alias string) *Users {
	cm := gql.NewColumnMaker("User", "users").As(alias)
	t := &Users{
		Table: gql.NewTableHelper("users", alias),

		ID:   cm.Col("ID", "id"),
		Name: cm.Col("Name", "name"),
	}
	t.ModelCollectorMaker = cllct.NewModelCollectorMaker(t.Columns(), alias)
	return t
}

func (t *Users) As(alias string) *Users { return NewUsers(alias) }
func (t *Users) All() gql.ExprListExpr  { return gql.AllCols(t.Columns()) }
func (t *Users) Columns() []gql.Column {
	return []gql.Column{t.ID, t.Name}
}

type Cities struct {
	gql.Table
	*cllct.ModelCollectorMaker

	ID           gql.Column
	Name         gql.Column
	PrefectureID gql.Column
}

func NewCities(alias string) *Cities {
	cm := gql.NewColumnMaker("City", "cities").As(alias)
	t := &Cities{
		Table: gql.NewTableHelper("cities", alias),

		ID:           cm.Col("ID", "id"),
		Name:         cm.Col("Name", "name"),
		PrefectureID: cm.Col("PrefectureID", "prefecture_id"),
	}
	t.ModelCollectorMaker = cllct.NewModelCollectorMaker(t.Columns(), alias)
	return t
}

func (t *Cities) As(alias string) *Cities { return NewCities(alias) }
func (t *Cities) All() gql.ExprListExpr   { return gql.AllCols(t.Columns()) }
func (t *Cities) Columns() []gql.Column {
	return []gql.Column{t.ID, t.Name, t.PrefectureID}
}

type Builder struct {
	*goq.Builder

	Users  *Users
	Cities *Cities
}

func NewBuilder(dl dialect.Dialect) *Builder {
	return &Builder{
		Builder: goq.NewBuilder(dl),

		Users:  NewUsers(""),
		Cities: NewCities(""),
	}
}
