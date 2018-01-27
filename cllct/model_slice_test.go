package cllct_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/ryym/goq/cllct"
	"github.com/ryym/goq/gql"
)

// Model collectors collect results regardless of
// selected column orders and their aliases.
func TestModelSliceCollector(t *testing.T) {
	rows := [][]interface{}{
		{"unrelated", "bob", 250},
		{"unrelated", "alice", 101},
	}

	selects := []gql.Selection{
		sel("", "foo", "Bar"),
		sel("", "User", "Name"),
		sel("", "User", "ID"),
	}

	users := NewUsers("")

	var got []User
	execCollector([]cllct.Collector{
		users.ToSlice(&got),
	}, rows, selects, nil)

	want := []User{
		{ID: 250, Name: "bob"},
		{ID: 101, Name: "alice"},
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}
