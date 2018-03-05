package cllct_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ryym/goq/cllct"
)

func TestRowMapCollector(t *testing.T) {
	rows := [][]interface{}{
		{1, "foo", true},
	}

	cl := cllct.NewMaker()

	var got map[string]interface{}
	execCollector([]cllct.Collector{
		cl.ToRowMap(&got),
	}, rows, nil, []string{"a", "b", "c"})

	want := map[string]interface{}{"a": 1, "b": "foo", "c": true}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Error(diff)
	}
}
