package cllct_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/ryym/goq/cllct"
	"github.com/ryym/goq/goql"
)

func TestModelMapCollector(t *testing.T) {
	rows := [][]interface{}{
		{"_", "bob", 250},
		{"_", "alice", 101},
		{"_", "carol", 18},
	}

	selects := []goql.Selection{
		sel("", "foo", "Bar"),
		sel("", "User", "Name"),
		sel("", "User", "ID"),
	}

	users := NewUsers("")

	var got map[int]User
	execCollector([]cllct.Collector{
		users.ToMap(&got),
	}, rows, selects, nil)

	want := map[int]User{
		250: User{ID: 250, Name: "bob"},
		101: User{ID: 101, Name: "alice"},
		18:  User{ID: 18, Name: "carol"},
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}
