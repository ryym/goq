package cllct_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/ryym/goq/cllct"
	"github.com/ryym/goq/gql"
)

func TestModelUniqSliceCollector(t *testing.T) {
	users := NewUsers("")
	maker := cllct.NewModelCollectorMaker(users.Columns(), "")

	var got []User
	cl := maker.ToUniqSlice(&got)

	rows := [][]interface{}{
		{2, "_", "bob"},
		{1, "_", "alice"},
		{2, "_", "bob"},
		{1, "_", "alice"},
		{1, "_", "alice"},
		{3, "_", "carol"},
		{2, "_", "bob"},
		{2, "_", "bob"},
		{3, "_", "carol"},
		{2, "_", "bob"},
	}

	selects := []gql.Selection{
		sel("", "User", "ID"),
		sel("", "foo", "Unrelated"),
		sel("", "User", "Name"),
	}

	execCollector(cl, rows, selects, nil)

	want := []User{
		{ID: 2, Name: "bob"},
		{ID: 1, Name: "alice"},
		{ID: 3, Name: "carol"},
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}
