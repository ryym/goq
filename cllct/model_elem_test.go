package cllct_test

import (
	"reflect"
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

	row := []interface{}{"alice", 101}

	selects := []gql.Selection{
		sel("", "User", "Name"),
		sel("", "User", "ID"),
	}

	cl.Init(selects, []string{"a", "b"})

	ptrs := make([]interface{}, len(selects))
	cl.Next(ptrs)
	for i, p := range ptrs {
		if p != nil {
			reflect.ValueOf(p).Elem().Set(reflect.ValueOf(row[i]))
		}
	}
	cl.AfterScan(ptrs)

	want := User{ID: 101, Name: "alice"}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}
