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

	ID   *gql.Column
	Name *gql.Column
}

func NewCountries(alias string) *Countries {
	cm := gql.NewColumnMaker("Country", "countries").As(alias)
	t := &Countries{
		Table: gql.NewTable("countries", alias),

		ID:   cm.Col("ID", "id").PK().Bld(),
		Name: cm.Col("Name", "name").Bld(),
	}
	t.ModelCollectorMaker = cllct.NewModelCollectorMaker(t.Columns(), alias)
	return t
}

func (t *Countries) As(alias string) *Countries { return NewCountries(alias) }
func (t *Countries) All() gql.ExprListExpr      { return gql.AllCols(t.Columns()) }
func (t *Countries) Columns() []*gql.Column {
	return []*gql.Column{t.ID, t.Name}
}

type Cities struct {
	gql.Table
	*cllct.ModelCollectorMaker

	ID        *gql.Column
	Name      *gql.Column
	CountryID *gql.Column
}

func NewCities(alias string) *Cities {
	cm := gql.NewColumnMaker("City", "cities").As(alias)
	t := &Cities{
		Table: gql.NewTable("cities", alias),

		ID:        cm.Col("ID", "id").PK().Bld(),
		Name:      cm.Col("Name", "name").Bld(),
		CountryID: cm.Col("CountryID", "country_id").Bld(),
	}
	t.ModelCollectorMaker = cllct.NewModelCollectorMaker(t.Columns(), alias)
	return t
}

func (t *Cities) As(alias string) *Cities { return NewCities(alias) }
func (t *Cities) All() gql.ExprListExpr   { return gql.AllCols(t.Columns()) }
func (t *Cities) Columns() []*gql.Column {
	return []*gql.Column{t.ID, t.Name, t.CountryID}
}

type Addresses struct {
	gql.Table
	*cllct.ModelCollectorMaker

	ID     *gql.Column
	Name   *gql.Column
	CityID *gql.Column
}

func NewAddresses(alias string) *Addresses {
	cm := gql.NewColumnMaker("Address", "addresses").As(alias)
	t := &Addresses{
		Table: gql.NewTable("addresses", alias),

		ID:     cm.Col("ID", "id").PK().Bld(),
		Name:   cm.Col("Name", "name").Bld(),
		CityID: cm.Col("CityID", "city_id").Bld(),
	}
	t.ModelCollectorMaker = cllct.NewModelCollectorMaker(t.Columns(), alias)
	return t
}

func (t *Addresses) As(alias string) *Addresses { return NewAddresses(alias) }
func (t *Addresses) All() gql.ExprListExpr      { return gql.AllCols(t.Columns()) }
func (t *Addresses) Columns() []*gql.Column {
	return []*gql.Column{t.ID, t.Name, t.CityID}
}

type Techs struct {
	gql.Table
	*cllct.ModelCollectorMaker

	ID   *gql.Column
	Name *gql.Column
	Desc *gql.Column
}

func NewTechs(alias string) *Techs {
	cm := gql.NewColumnMaker("Tech", "technologies").As(alias)
	t := &Techs{
		Table: gql.NewTable("technologies", alias),

		ID:   cm.Col("ID", "id").PK().Bld(),
		Name: cm.Col("Name", "name").Bld(),
		Desc: cm.Col("Desc", "description").Bld(),
	}
	t.ModelCollectorMaker = cllct.NewModelCollectorMaker(t.Columns(), alias)
	return t
}

func (t *Techs) As(alias string) *Techs { return NewTechs(alias) }
func (t *Techs) All() gql.ExprListExpr  { return gql.AllCols(t.Columns()) }
func (t *Techs) Columns() []*gql.Column {
	return []*gql.Column{t.ID, t.Name, t.Desc}
}

type Builder struct {
	*goq.Builder

	Countries *Countries
	Cities    *Cities
	Addresses *Addresses
	Techs     *Techs
}

func NewBuilder(dl dialect.Dialect) *Builder {
	return &Builder{
		Builder: goq.NewBuilder(dl),

		Countries: NewCountries(""),
		Cities:    NewCities(""),
		Addresses: NewAddresses(""),
		Techs:     NewTechs(""),
	}
}
