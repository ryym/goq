package cllct_test

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
	. "github.com/ryym/goq/cllct"
	"github.com/ryym/goq/gql"
)

func TestRowMapSliceCollector(t *testing.T) {
	var got []map[string]interface{}
	cl := NewRowMapSliceCollector(&got)

	rows := [][]interface{}{
		{1, "foo", true},
		{2, "bar", false},
	}

	cl.Init(
		make([]gql.Selection, 3),
		[]string{"a", "b", "c"},
	)

	for _, row := range rows {
		ptrs := make([]interface{}, 3)
		cl.Next(ptrs)
		for i, p := range ptrs {
			reflect.ValueOf(p).Elem().Set(reflect.ValueOf(row[i]))
		}
		cl.AfterScan(ptrs)
	}

	want := []map[string]interface{}{
		{"a": 1, "b": "foo", "c": true},
		{"a": 2, "b": "bar", "c": false},
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}
