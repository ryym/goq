package sample

import (
	"github.com/ryym/goq"
	"github.com/ryym/goq/_sample/models"
	"github.com/ryym/goq/dialect"
	"github.com/ryym/goq/gql"
)

type Users struct {
	gql.TableHelper
	model models.User

	ID   gql.Column
	Name gql.Column
}

func NewUsers() *Users {
	cm := gql.NewColumnMaker("User", "users")
	return &Users{
		TableHelper: gql.NewTableHelper("users", ""),
		model:       models.User{},

		ID:   cm.Col("ID", "id"),
		Name: cm.Col("Name", "name"),
	}
}

func (t *Users) All() gql.ExprListExpr { return gql.AllCols(t.Columns()) }
func (t *Users) Columns() []gql.Column {
	return []gql.Column{t.ID, t.Name}
}
func (t *Users) As(alias string) *Users {
	t2 := *t
	t2.TableHelper = gql.NewTableHelper(t.TableHelper.TableName(), alias)
	gql.CopyTableAs(alias, t, &t2)
	return &t2
}

type Cities struct {
	gql.TableHelper
	model City

	ID           gql.Column
	Name         gql.Column
	PrefectureID gql.Column
}

func NewCities() *Cities {
	cm := gql.NewColumnMaker("City", "cities")
	return &Cities{
		TableHelper: gql.NewTableHelper("cities", ""),
		model:       City{},

		ID:           cm.Col("ID", "id"),
		Name:         cm.Col("Name", "name"),
		PrefectureID: cm.Col("PrefectureID", "prefecture_id"),
	}
}

func (t *Cities) All() gql.ExprListExpr { return gql.AllCols(t.Columns()) }
func (t *Cities) Columns() []gql.Column {
	return []gql.Column{t.ID, t.Name, t.PrefectureID}
}
func (t *Cities) As(alias string) *Cities {
	t2 := *t
	t2.TableHelper = gql.NewTableHelper(t.TableHelper.TableName(), alias)
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
