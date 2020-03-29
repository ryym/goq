package goq_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ryym/goq"
)

func TestModelUniqSliceCollector(t *testing.T) {
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

	selects := []goq.Selection{
		sel("", "User", "ID"),
		sel("", "foo", "Unrelated"),
		sel("", "User", "Name"),
	}

	users := NewUsers("")

	var got []User
	execCollector([]goq.Collector{
		users.ToUniqSlice(&got),
	}, rows, selects, nil)

	want := []User{
		{ID: 2, Name: "bob"},
		{ID: 1, Name: "alice"},
		{ID: 3, Name: "carol"},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Error(diff)
	}
}
