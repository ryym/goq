package goq_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ryym/goq"
	"github.com/ryym/goq/dialect"
)

func TestSliceMapCollector(t *testing.T) {
	rows := [][]interface{}{
		{1, 1, "a"},
		{1, 2, "b"},
		{1, 3, "c"},
		{1, 4, "d"},
		{2, 5, "e"},
		{2, 6, "f"},
		{2, 7, "g"},
	}

	q := goq.NewQueryBuilder(dialect.Generic())
	cl := goq.NewMaker()
	names := []string{"country_id", "id", "name"}

	var got map[int][]City
	execCollector([]goq.Collector{
		cl.ToSliceMap(&got).By(q.Name("country_id")),
	}, rows, nil, names)

	want := map[int][]City{
		1: []City{
			{ID: 1, Name: "a", CountryID: 1},
			{ID: 2, Name: "b", CountryID: 1},
			{ID: 3, Name: "c", CountryID: 1},
			{ID: 4, Name: "d", CountryID: 1},
		},
		2: []City{
			{ID: 5, Name: "e", CountryID: 2},
			{ID: 6, Name: "f", CountryID: 2},
			{ID: 7, Name: "g", CountryID: 2},
		},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Error(diff)
	}
}
