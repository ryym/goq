package cllct_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/ryym/goq/cllct"
	"github.com/ryym/goq/dialect"
	"github.com/ryym/goq/gql"
)

func TestMapCollector(t *testing.T) {
	rows := [][]interface{}{
		{3, 1, "foo"},
		{4, 53, "bar"},
	}

	q := gql.NewBuilder(dialect.Generic())
	cl := cllct.NewMaker()
	names := []string{"id", "country_id", "name"}

	var got map[int]City
	execCollector([]cllct.Collector{
		cl.ToMap(&got).By(q.Name("country_id")),
	}, rows, nil, names)

	want := map[int]City{
		1:  {ID: 3, Name: "foo", CountryID: 1},
		53: {ID: 4, Name: "bar", CountryID: 53},
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}
