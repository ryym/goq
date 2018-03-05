package cllct_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ryym/goq/cllct"
)

func TestElemCollector(t *testing.T) {
	rows := [][]interface{}{
		{3, 1, "foo"},
	}

	cl := cllct.NewMaker()
	names := []string{"id", "country_id", "name"}

	var got City
	execCollector([]cllct.Collector{
		cl.ToElem(&got),
	}, rows, nil, names)

	want := City{ID: 3, Name: "foo", CountryID: 1}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Error(diff)
	}
}
