package cllct_test

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
	. "github.com/ryym/goq/cllct"
	"github.com/ryym/goq/gql"
)

func sel(alias, strct, field string) gql.Selection {
	return gql.Selection{TableAlias: alias, StructName: strct, FieldName: field}
}

// Model collectors collect results regardless of
// selected column orders and their aliases.
func TestModelSliceCollector(t *testing.T) {
	users := NewUsers("")
	maker := NewModelCollectorMaker(users.Columns(), "")

	var got []User
	cl := maker.ToSlice(&got)

	rows := [][]interface{}{
		{"unrelated", "bob", 250},
		{"unrelated", "alice", 101},
	}

	selects := []gql.Selection{
		sel("", "foo", "Bar"),
		sel("", "User", "Name"),
		sel("", "User", "ID"),
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

	want := []User{
		{ID: 250, Name: "bob"},
		{ID: 101, Name: "alice"},
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}
