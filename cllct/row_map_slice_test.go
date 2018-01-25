package cllct_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/ryym/goq/cllct"
)

func TestRowMapSliceCollector(t *testing.T) {
	var got []map[string]interface{}
	cl := cllct.NewRowMapSliceCollector(&got)

	rows := [][]interface{}{
		{1, "foo", true},
		{2, "bar", false},
	}
	execCollector(cl, rows, nil, []string{"a", "b", "c"})

	want := []map[string]interface{}{
		{"a": 1, "b": "foo", "c": true},
		{"a": 2, "b": "bar", "c": false},
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}
