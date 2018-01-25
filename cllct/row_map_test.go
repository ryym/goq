package cllct_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/ryym/goq/cllct"
)

func TestRowMapCollector(t *testing.T) {
	var got map[string]interface{}
	cl := cllct.NewMaker().ToRowMap(&got)

	rows := [][]interface{}{
		{1, "foo", true},
	}
	execCollector(cl, rows, nil, []string{"a", "b", "c"})

	want := map[string]interface{}{"a": 1, "b": "foo", "c": true}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}
