package tests

import (
	"github.com/ryym/goq"
	"github.com/ryym/goq/cllct"
	"github.com/ryym/goq/dialect"
	"github.com/ryym/goq/gql"
)

type Countries struct {
	gql.Table
	*cllct.ModelCollectorMaker

	ID        *gql.Column
	Name      *gql.Column
	UpdatedAt *gql.Column
}

func NewCountries(alias string) *Countries {
	cm := gql.NewColumnMaker("Country", "countries").As(alias)
	t := &Countries{
		Table: gql.NewTable("countries", alias),

		ID:        cm.Col("ID", "id").Bld(),
		Name:      cm.Col("Name", "name").Bld(),
		UpdatedAt: cm.Col("UpdatedAt", "updated_at").Bld(),
	}
	t.ModelCollectorMaker = cllct.NewModelCollectorMaker(t.Columns(), alias)
	return t
}

func (t *Countries) As(alias string) *Countries { return NewCountries(alias) }
func (t *Countries) All() gql.ExprListExpr      { return gql.AllCols(t.Columns()) }
func (t *Countries) Columns() []*gql.Column {
	return []*gql.Column{t.ID, t.Name, t.UpdatedAt}
}

type Cities struct {
	gql.Table
	*cllct.ModelCollectorMaker

	ID        *gql.Column
	Name      *gql.Column
	CountryID *gql.Column
	UpdatedAt *gql.Column
}

func NewCities(alias string) *Cities {
	cm := gql.NewColumnMaker("City", "cities").As(alias)
	t := &Cities{
		Table: gql.NewTable("cities", alias),

		ID:        cm.Col("ID", "id").Bld(),
		Name:      cm.Col("Name", "name").Bld(),
		CountryID: cm.Col("CountryID", "country_id").Bld(),
		UpdatedAt: cm.Col("UpdatedAt", "updated_at").Bld(),
	}
	t.ModelCollectorMaker = cllct.NewModelCollectorMaker(t.Columns(), alias)
	return t
}

func (t *Cities) As(alias string) *Cities { return NewCities(alias) }
func (t *Cities) All() gql.ExprListExpr   { return gql.AllCols(t.Columns()) }
func (t *Cities) Columns() []*gql.Column {
	return []*gql.Column{t.ID, t.Name, t.CountryID, t.UpdatedAt}
}

type Builder struct {
	*goq.Builder

	Countries *Countries
	Cities    *Cities
}

func NewBuilder(dl dialect.Dialect) *Builder {
	return &Builder{
		Builder: goq.NewBuilder(dl),

		Countries: NewCountries(""),
		Cities:    NewCities(""),
	}
}
