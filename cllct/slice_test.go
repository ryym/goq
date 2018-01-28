package cllct_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/ryym/goq/cllct"
)

func TestSliceCollector(t *testing.T) {
	rows := [][]interface{}{
		{3, 1, "foo"},
		{4, 53, "bar"},
	}

	cl := cllct.NewMaker()
	names := []string{"id", "country_id", "name"}

	var got []City
	execCollector([]cllct.Collector{
		cl.ToSlice(&got),
	}, rows, nil, names)

	want := []City{
		{ID: 3, Name: "foo", CountryID: 1},
		{ID: 4, Name: "bar", CountryID: 53},
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}
