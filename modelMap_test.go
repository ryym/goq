package goq_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ryym/goq"
)

func TestModelMapCollector(t *testing.T) {
	rows := [][]interface{}{
		{"_", "bob", 250},
		{"_", "alice", 101},
		{"_", "carol", 18},
	}

	selects := []goq.Selection{
		sel("", "foo", "Bar"),
		sel("", "User", "Name"),
		sel("", "User", "ID"),
	}

	users := NewUsers("")

	var got map[int]User
	execCollector([]goq.Collector{
		users.ToMap(&got),
	}, rows, selects, nil)

	want := map[int]User{
		250: User{ID: 250, Name: "bob"},
		101: User{ID: 101, Name: "alice"},
		18:  User{ID: 18, Name: "carol"},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Error(diff)
	}
}
