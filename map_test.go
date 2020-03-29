package goq_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ryym/goq"
	"github.com/ryym/goq/dialect"
	"github.com/ryym/goq/goql"
)

func TestMapCollector(t *testing.T) {
	rows := [][]interface{}{
		{3, 1, "foo"},
		{4, 53, "bar"},
	}

	q := goql.NewQueryBuilder(dialect.Generic())
	cl := goq.NewMaker()
	names := []string{"id", "country_id", "name"}

	var got map[int]City
	execCollector([]goq.Collector{
		cl.ToMap(&got).By(q.Name("country_id")),
	}, rows, nil, names)

	want := map[int]City{
		1:  {ID: 3, Name: "foo", CountryID: 1},
		53: {ID: 4, Name: "bar", CountryID: 53},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Error(diff)
	}
}

func TestInvalidMapCollector(t *testing.T) {
	cl := goq.NewMaker()
	q := goql.NewQueryBuilder(dialect.Generic())
	initConf := goq.NewInitConf([]goql.Selection{}, []string{})
	var err error
	var cllctor *goq.MapCollector

	// Not a pointer
	cllctor = cl.ToMap(map[int]interface{}{}).By(q.Name("id"))
	_, err = goq.InitCollectors([]goq.Collector{cllctor}, initConf)
	if err == nil {
		t.Error("ToMap accepts not a pointer")
	}

	// Invalid pointer
	var slice []interface{}
	cllctor = cl.ToMap(&slice).By(q.Name("id"))
	_, err = goq.InitCollectors([]goq.Collector{cllctor}, initConf)
	if err == nil {
		t.Error("ToMap accepts a pointer not to map")
	}

	// Map key not collected
	var mp map[int]struct{}
	cllctor = cl.ToMap(&mp).By(q.Name("id"))
	_, err = goq.InitCollectors([]goq.Collector{cllctor}, initConf)
	if err == nil {
		t.Error("ToMap accespts invalid map key")
	}
}
