// DO NOT EDIT. This code is generated by Goq.
// https://github.com/ryym/goq

package tests

import (
	"github.com/ryym/goq"
	"github.com/ryym/goq/dialect"
)

type Countries struct {
	goq.Table
	*goq.ModelCollectorMaker

	ID   *goq.Column
	Name *goq.Column
}

func NewCountries(alias string) *Countries {
	cm := goq.NewColumnMaker("Country", "countries").As(alias)
	t := &Countries{

		ID:   cm.Col("ID", "id").PK().Bld(),
		Name: cm.Col("Name", "name").Bld(),
	}
	cols := []*goq.Column{t.ID, t.Name}
	t.Table = goq.NewTable("countries", alias, cols)
	t.ModelCollectorMaker = goq.NewModelCollectorMaker(cols, alias)
	return t
}

func (t *Countries) As(alias string) *Countries { return NewCountries(alias) }

type Cities struct {
	goq.Table
	*goq.ModelCollectorMaker

	ID        *goq.Column
	Name      *goq.Column
	CountryID *goq.Column
}

func NewCities(alias string) *Cities {
	cm := goq.NewColumnMaker("City", "cities").As(alias)
	t := &Cities{

		ID:        cm.Col("ID", "id").PK().Bld(),
		Name:      cm.Col("Name", "name").Bld(),
		CountryID: cm.Col("CountryID", "country_id").Bld(),
	}
	cols := []*goq.Column{t.ID, t.Name, t.CountryID}
	t.Table = goq.NewTable("cities", alias, cols)
	t.ModelCollectorMaker = goq.NewModelCollectorMaker(cols, alias)
	return t
}

func (t *Cities) As(alias string) *Cities { return NewCities(alias) }

type Addresses struct {
	goq.Table
	*goq.ModelCollectorMaker

	ID     *goq.Column
	Name   *goq.Column
	CityID *goq.Column
}

func NewAddresses(alias string) *Addresses {
	cm := goq.NewColumnMaker("Address", "addresses").As(alias)
	t := &Addresses{

		ID:     cm.Col("ID", "id").PK().Bld(),
		Name:   cm.Col("Name", "name").Bld(),
		CityID: cm.Col("CityID", "city_id").Bld(),
	}
	cols := []*goq.Column{t.ID, t.Name, t.CityID}
	t.Table = goq.NewTable("addresses", alias, cols)
	t.ModelCollectorMaker = goq.NewModelCollectorMaker(cols, alias)
	return t
}

func (t *Addresses) As(alias string) *Addresses { return NewAddresses(alias) }

type Techs struct {
	goq.Table
	*goq.ModelCollectorMaker

	ID   *goq.Column
	Name *goq.Column
	Desc *goq.Column
}

func NewTechs(alias string) *Techs {
	cm := goq.NewColumnMaker("Tech", "technologies").As(alias)
	t := &Techs{

		ID:   cm.Col("ID", "id").PK().Bld(),
		Name: cm.Col("Name", "name").Bld(),
		Desc: cm.Col("Desc", "description").Bld(),
	}
	cols := []*goq.Column{t.ID, t.Name, t.Desc}
	t.Table = goq.NewTable("technologies", alias, cols)
	t.ModelCollectorMaker = goq.NewModelCollectorMaker(cols, alias)
	return t
}

func (t *Techs) As(alias string) *Techs { return NewTechs(alias) }

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
