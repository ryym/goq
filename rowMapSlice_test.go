package goq_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ryym/goq"
)

func TestRowMapSliceCollector(t *testing.T) {
	rows := [][]interface{}{
		{1, "foo", true},
		{2, "bar", false},
	}

	cl := goq.NewMaker()

	var got []map[string]interface{}
	execCollector([]goq.Collector{
		cl.ToRowMapSlice(&got),
	}, rows, nil, []string{"a", "b", "c"})

	want := []map[string]interface{}{
		{"a": 1, "b": "foo", "c": true},
		{"a": 2, "b": "bar", "c": false},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Error(diff)
	}
}
