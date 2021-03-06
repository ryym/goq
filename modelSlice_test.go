package goq_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ryym/goq"
)

// Model collectors collect results regardless of
// selected column orders and their aliases.
func TestModelSliceCollector(t *testing.T) {
	rows := [][]interface{}{
		{"unrelated", "bob", 250},
		{"unrelated", "alice", 101},
	}

	selects := []goq.Selection{
		sel("", "foo", "Bar"),
		sel("", "User", "Name"),
		sel("", "User", "ID"),
	}

	users := NewUsers("")

	var got []User
	execCollector([]goq.Collector{
		users.ToSlice(&got),
	}, rows, selects, nil)

	want := []User{
		{ID: 250, Name: "bob"},
		{ID: 101, Name: "alice"},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Error(diff)
	}
}
