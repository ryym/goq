package goq_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ryym/goq"
)

func TestSliceCollector(t *testing.T) {
	rows := [][]interface{}{
		{3, 1, "foo"},
		{4, 53, "bar"},
	}

	cl := goq.NewMaker()
	names := []string{"id", "country_id", "name"}

	var got []City
	execCollector([]goq.Collector{
		cl.ToSlice(&got),
	}, rows, nil, names)

	want := []City{
		{ID: 3, Name: "foo", CountryID: 1},
		{ID: 4, Name: "bar", CountryID: 53},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Error(diff)
	}
}
