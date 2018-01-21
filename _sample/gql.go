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

func NewUsers() *Users {
	cm := gql.NewColumnMaker("User", "users")
	t := &Users{
		Table: gql.NewTableHelper("users", ""),

		ID:   cm.Col("ID", "id"),
		Name: cm.Col("Name", "name"),
	}
	t.ModelCollectorMaker = cllct.NewModelCollectorMaker(t.Columns(), "")
	return t
}

func (t *Users) All() gql.ExprListExpr { return gql.AllCols(t.Columns()) }
func (t *Users) Columns() []gql.Column {
	return []gql.Column{t.ID, t.Name}
}
func (t *Users) As(alias string) *Users {
	t2 := *t
	t2.Table = gql.NewTableHelper(t.Table.TableName(), alias)
	t2.ModelCollectorMaker = cllct.NewModelCollectorMaker(t.Columns(), alias)
	gql.CopyTableAs(alias, t, &t2)
	return &t2
}

type Cities struct {
	gql.Table
	*cllct.ModelCollectorMaker

	ID           gql.Column
	Name         gql.Column
	PrefectureID gql.Column
}

func NewCities() *Cities {
	cm := gql.NewColumnMaker("City", "cities")
	t := &Cities{
		Table: gql.NewTableHelper("cities", ""),

		ID:           cm.Col("ID", "id"),
		Name:         cm.Col("Name", "name"),
		PrefectureID: cm.Col("PrefectureID", "prefecture_id"),
	}
	t.ModelCollectorMaker = cllct.NewModelCollectorMaker(t.Columns(), "")
	return t
}

func (t *Cities) All() gql.ExprListExpr { return gql.AllCols(t.Columns()) }
func (t *Cities) Columns() []gql.Column {
	return []gql.Column{t.ID, t.Name, t.PrefectureID}
}
func (t *Cities) As(alias string) *Cities {
	t2 := *t
	t2.Table = gql.NewTableHelper(t.Table.TableName(), alias)
	t2.ModelCollectorMaker = cllct.NewModelCollectorMaker(t.Columns(), alias)
	gql.CopyTableAs(alias, t, &t2)
	return &t2
}

type Builder struct {
	*goq.Builder

	Users  *Users
	Cities *Cities
}

func NewBuilder(dl dialect.Dialect) *Builder {
	return &Builder{
		Builder: goq.NewBuilder(dl),

		Users:  NewUsers(),
		Cities: NewCities(),
	}
}
