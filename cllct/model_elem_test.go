package cllct_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/ryym/goq/cllct"
	"github.com/ryym/goq/gql"
)

func TestModelElemCollector(t *testing.T) {
	users := NewUsers("")
	maker := cllct.NewModelCollectorMaker(users.Columns(), "")

	var got User
	cl := maker.ToElem(&got)

	rows := [][]interface{}{
		{"alice", 101},
	}

	selects := []gql.Selection{
		sel("", "User", "Name"),
		sel("", "User", "ID"),
	}

	execCollector(cl, rows, selects, nil)

	want := User{ID: 101, Name: "alice"}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}
