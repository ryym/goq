package cllct_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/ryym/goq/cllct"
	"github.com/ryym/goq/goql"
)

func TestModelElemCollector(t *testing.T) {
	rows := [][]interface{}{
		{"alice", 101},
	}

	selects := []goql.Selection{
		sel("", "User", "Name"),
		sel("", "User", "ID"),
	}

	users := NewUsers("")

	var got User
	execCollector([]cllct.Collector{
		users.ToElem(&got),
	}, rows, selects, nil)

	want := User{ID: 101, Name: "alice"}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}
