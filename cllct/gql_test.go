// DO NOT EDIT. This code is generated by Goq.
// https://github.com/ryym/goq

package cllct_test

import (
	"github.com/ryym/goq"
	"github.com/ryym/goq/cllct"
	"github.com/ryym/goq/dialect"
	"github.com/ryym/goq/gql"
)

type Users struct {
	gql.Table
	*cllct.ModelCollectorMaker

	ID   *gql.Column
	Name *gql.Column
}

func NewUsers(alias string) *Users {
	cm := gql.NewColumnMaker("User", "users").As(alias)
	t := &Users{
		Table: gql.NewTable("users", alias),

		ID:   cm.Col("ID", "id").PK().Bld(),
		Name: cm.Col("Name", "name").Bld(),
	}
	t.ModelCollectorMaker = cllct.NewModelCollectorMaker(t.Columns(), alias)
	return t
}

func (t *Users) As(alias string) *Users   { return NewUsers(alias) }
func (t *Users) All() *gql.ColumnListExpr { return gql.AllCols(t.Columns()) }
func (t *Users) Columns() []*gql.Column {
	return []*gql.Column{t.ID, t.Name}
}

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
func (t *Countries) All() *gql.ColumnListExpr   { return gql.AllCols(t.Columns()) }
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

func (t *Cities) As(alias string) *Cities  { return NewCities(alias) }
func (t *Cities) All() *gql.ColumnListExpr { return gql.AllCols(t.Columns()) }
func (t *Cities) Columns() []*gql.Column {
	return []*gql.Column{t.ID, t.Name, t.CountryID}
}

type Addresses struct {
	gql.Table
	*cllct.ModelCollectorMaker

	ID         *gql.Column
	Address1   *gql.Column
	Address2   *gql.Column
	CityID     *gql.Column
	PostalCode *gql.Column
}

func NewAddresses(alias string) *Addresses {
	cm := gql.NewColumnMaker("Address", "addresses").As(alias)
	t := &Addresses{
		Table: gql.NewTable("addresses", alias),

		ID:         cm.Col("ID", "id").PK().Bld(),
		Address1:   cm.Col("Address1", "address1").Bld(),
		Address2:   cm.Col("Address2", "address2").Bld(),
		CityID:     cm.Col("CityID", "city_id").Bld(),
		PostalCode: cm.Col("PostalCode", "postal_code").Bld(),
	}
	t.ModelCollectorMaker = cllct.NewModelCollectorMaker(t.Columns(), alias)
	return t
}

func (t *Addresses) As(alias string) *Addresses { return NewAddresses(alias) }
func (t *Addresses) All() *gql.ColumnListExpr   { return gql.AllCols(t.Columns()) }
func (t *Addresses) Columns() []*gql.Column {
	return []*gql.Column{t.ID, t.Address1, t.Address2, t.CityID, t.PostalCode}
}

type Builder struct {
	*goq.Builder

	Users     *Users
	Countries *Countries
	Cities    *Cities
	Addresses *Addresses
}

func NewBuilder(dl dialect.Dialect) *Builder {
	return &Builder{
		Builder: goq.NewBuilder(dl),

		Users:     NewUsers(""),
		Countries: NewCountries(""),
		Cities:    NewCities(""),
		Addresses: NewAddresses(""),
	}
}
