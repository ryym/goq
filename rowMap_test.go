package goq_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ryym/goq"
)

func TestRowMapCollector(t *testing.T) {
	rows := [][]interface{}{
		{1, "foo", true},
	}

	cl := goq.NewMaker()

	var got map[string]interface{}
	execCollector([]goq.Collector{
		cl.ToRowMap(&got),
	}, rows, nil, []string{"a", "b", "c"})

	want := map[string]interface{}{"a": 1, "b": "foo", "c": true}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Error(diff)
	}
}
