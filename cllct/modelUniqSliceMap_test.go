package cllct_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/ryym/goq/cllct"
	"github.com/ryym/goq/gql"
)

func TestModelUniqSliceMapCollector(t *testing.T) {
	rows := [][]interface{}{
		{8, "japan", 1, "tokyo"},
		{8, "japan", 1, "tokyo"},
		{8, "japan", 1, "tokyo"},
		{8, "japan", 2, "oosaka"},
		{8, "japan", 2, "oosaka"},
		{8, "japan", 3, "hiroshima"},
		{3, "us", 4, "newyork"},
		{3, "us", 4, "newyork"},
		{3, "us", 4, "newyork"},
		{3, "us", 5, "losangeles"},
		{3, "us", 5, "losangeles"},
		{3, "us", 5, "losangeles"},
		{3, "us", 5, "losangeles"},
		{3, "us", 6, "chicago"},
		{3, "us", 7, "houston"},
		{3, "us", 7, "houston"},
	}

	selects := []gql.Selection{
		sel("", "Country", "ID"),
		sel("", "Country", "Name"),
		sel("", "City", "ID"),
		sel("", "City", "Name"),
	}

	countries := NewCountries("")
	cities := NewCities("")

	var countryID int
	var got map[int][]City
	err := execCollector([]cllct.Collector{
		cities.ToUniqSliceMap(&got).ByWith(&countryID, countries.ID),
	}, rows, selects, nil)
	if err != nil {
		t.Fatal(err)
	}

	want := map[int][]City{
		8: []City{
			{ID: 1, Name: "tokyo"},
			{ID: 2, Name: "oosaka"},
			{ID: 3, Name: "hiroshima"},
		},
		3: []City{
			{ID: 4, Name: "newyork"},
			{ID: 5, Name: "losangeles"},
			{ID: 6, Name: "chicago"},
			{ID: 7, Name: "houston"},
		},
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}
