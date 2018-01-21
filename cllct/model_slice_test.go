package cllct

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/ryym/goq/gql"
)

type user struct {
	ID   int
	Name string
}

type Users struct {
	gql.TableHelper
	ID   gql.Column
	Name gql.Column
}

func NewUsers() *Users {
	cm := gql.NewColumnMaker("user", "users")
	return &Users{
		TableHelper: gql.NewTableHelper("users", ""),
		ID:          cm.Col("ID", "id"),
		Name:        cm.Col("Name", "name"),
	}
}

func (t *Users) Columns() []gql.Column {
	return []gql.Column{t.ID, t.Name}
}

func sel(alias, strct, field string) gql.Selection {
	return gql.Selection{TableAlias: alias, StructName: strct, FieldName: field}
}

// Model collectors collect results regardless of
// selected column orders and their aliases.
func TestModelSliceCollector(t *testing.T) {
	users := NewUsers()
	maker := NewModelCollectorMaker(user{}, users.Columns(), "")

	var got []user
	cl := maker.ToSlice(&got)

	rows := [][]interface{}{
		{"unrelated", "bob", 250},
		{"unrelated", "alice", 101},
	}

	selects := []gql.Selection{
		sel("", "foo", "Bar"),
		sel("", "user", "Name"),
		sel("", "user", "ID"),
	}

	cl.Init(selects, []string{"a", "b", "c"})

	for _, row := range rows {
		ptrs := make([]interface{}, len(selects))
		cl.Next(ptrs)
		for i, p := range ptrs {
			if p != nil {
				reflect.ValueOf(p).Elem().Set(reflect.ValueOf(row[i]))
			}
		}
		cl.AfterScan(ptrs)
	}

	want := []user{
		{ID: 250, Name: "bob"},
		{ID: 101, Name: "alice"},
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}
