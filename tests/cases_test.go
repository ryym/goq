package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/ryym/goq"
)

var defaultUpdatedAt = time.Date(2000, time.January, 1, 9, 0, 0, 0, time.UTC)

var testCases = []testCase{
	{
		name: "select first one from multiple records",
		data: `
			INSERT INTO cities VALUES (1, 'tokyo', 87, '2006-02-15 09:45:25');
			INSERT INTO cities VALUES (2, 'hokkaido', 87, '2006-02-15 09:45:25');
			INSERT INTO cities VALUES (3, 'okinawa', 87, '2006-02-15 09:45:25');
		`,
		run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
			q := z.Select(z.Cities.All()).From(z.Cities).OrderBy(z.Cities.Name)
			var city City
			err := tx.Query(q).First(z.Cities.ToElem(&city))
			if err != nil {
				return err
			}

			want := City{
				ID:        2,
				Name:      "hokkaido",
				CountryID: 87,
				UpdatedAt: time.Date(2006, time.February, 15, 9, 45, 25, 0, time.UTC),
			}
			if diff := deep.Equal(city, want); diff != nil {
				t.Log(q.Construct())
				return fmt.Errorf("%s", diff)
			}
			return nil
		},
	},
	{
		name: "select multiple records as slice",
		data: `
			INSERT INTO cities (id, name, country_id) VALUES (10, 'newyork', 5);
			INSERT INTO cities (id, name, country_id) VALUES (12, 'chicago', 5);
			INSERT INTO cities (id, name, country_id) VALUES (14, 'seattle', 5);
		`,
		run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
			q := z.Select(
				z.Cities.ID, z.Cities.Name,
			).From(z.Cities).OrderBy(z.Cities.ID.Desc())

			var cities []City
			err := tx.Query(q).Collect(z.Cities.ToSlice(&cities))
			if err != nil {
				return err
			}

			want := []City{
				{ID: 14, Name: "seattle"},
				{ID: 12, Name: "chicago"},
				{ID: 10, Name: "newyork"},
			}
			if diff := deep.Equal(cities, want); diff != nil {
				t.Log(q.Construct())
				return fmt.Errorf("%s", diff)
			}
			return nil
		},
	},
	{
		name: "select multiple records as map",
		data: `
			INSERT INTO cities (id, name, country_id) VALUES (10, 'newyork', 5);
			INSERT INTO cities (id, name, country_id) VALUES (12, 'chicago', 5);
			INSERT INTO cities (id, name, country_id) VALUES (14, 'seattle', 5);
		`,
		run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
			q := z.Select(
				z.Cities.ID, z.Cities.Name,
			).From(z.Cities).OrderBy(z.Cities.ID.Desc())

			var cities map[int]City
			err := tx.Query(q).Collect(z.Cities.ToMap(&cities))
			if err != nil {
				return err
			}

			want := map[int]City{
				10: {ID: 10, Name: "newyork"},
				12: {ID: 12, Name: "chicago"},
				14: {ID: 14, Name: "seattle"},
			}
			if diff := deep.Equal(cities, want); diff != nil {
				t.Log(q.Construct())
				return fmt.Errorf("%s", diff)
			}
			return nil
		},
	},
	{
		name: "select records by multiple slice collectors",
		data: `
			INSERT INTO countries (id, name) VALUES (8, 'Japan');
			INSERT INTO countries (id, name) VALUES (20, 'Nowhere');

			INSERT INTO cities (name, country_id) VALUES ('tokyo', 8);
			INSERT INTO cities (name, country_id) VALUES ('hokkaido', 8);
			INSERT INTO cities (name, country_id) VALUES ('okinawa', 8);
		`,
		run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
			q := z.Select(
				z.Countries.Name,
				z.Cities.Name,
				z.Cities.CountryID,
			).From(
				z.Countries,
			).Joins(
				z.LeftJoin(z.Cities).On(
					z.Cities.CountryID.Eq(z.Countries.ID),
				),
			).OrderBy(
				z.Countries.Name,
				z.Cities.Name,
			)

			var countries []Country
			var cities []City
			err := tx.Query(q).Collect(
				z.Countries.ToSlice(&countries),
				z.Cities.ToSlice(&cities),
			)
			if err != nil {
				return err
			}

			// Duplicated results.
			wantCountries := []Country{
				{Name: "Japan"},
				{Name: "Japan"},
				{Name: "Japan"},
				{Name: "Nowhere"},
			}
			if diff := deep.Equal(countries, wantCountries); diff != nil {
				t.Log(q.Construct())
				return fmt.Errorf("countries diff: %s", diff)
			}

			// Empty City exists because of LEFT JOIN.
			wantCities := []City{
				{Name: "hokkaido", CountryID: 8},
				{Name: "okinawa", CountryID: 8},
				{Name: "tokyo", CountryID: 8},
				{},
			}
			if diff := deep.Equal(cities, wantCities); diff != nil {
				t.Log(q.Construct())
				return fmt.Errorf("cities diff: %s", diff)
			}

			return nil
		},
	},
	{
		name: "select records uniquely by multiple slice collectors",
		data: `
			INSERT INTO countries (id, name) VALUES (8, 'Japan');
			INSERT INTO countries (id, name) VALUES (20, 'Somewhere');

			INSERT INTO cities (name, country_id) VALUES ('tokyo', 8);
			INSERT INTO cities (name, country_id) VALUES ('hokkaido', 8);
			INSERT INTO cities (name, country_id) VALUES ('okinawa', 8);
			INSERT INTO cities (name, country_id) VALUES ('foo', 20);
		`,
		run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
			q := z.Select(
				z.Countries.ID,
				z.Countries.Name,
				z.Cities.Name,
				z.Cities.CountryID,
			).From(
				z.Countries,
			).Joins(
				z.LeftJoin(z.Cities).On(
					z.Cities.CountryID.Eq(z.Countries.ID),
				),
			).OrderBy(
				z.Countries.Name,
				z.Cities.Name,
			)

			var countries []Country
			var cities []City
			err := tx.Query(q).Collect(
				z.Countries.ToUniqSlice(&countries),
				z.Cities.ToSlice(&cities),
			)
			if err != nil {
				return err
			}

			wantCountries := []Country{
				{ID: 8, Name: "Japan"},
				{ID: 20, Name: "Somewhere"},
			}
			if diff := deep.Equal(countries, wantCountries); diff != nil {
				t.Log(q.Construct())
				return fmt.Errorf("countries diff: %s", diff)
			}

			wantCities := []City{
				{Name: "hokkaido", CountryID: 8},
				{Name: "okinawa", CountryID: 8},
				{Name: "tokyo", CountryID: 8},
				{Name: "foo", CountryID: 20},
			}
			if diff := deep.Equal(cities, wantCities); diff != nil {
				t.Log(q.Construct())
				return fmt.Errorf("cities diff: %s", diff)
			}

			return nil
		},
	},
	{
		name: "select records by slice and map collectors",
		data: `
			INSERT INTO countries (id, name) VALUES (1, 'Japan');

			INSERT INTO cities (country_id, id, name) VALUES (1, 11, 'tokyo');
			INSERT INTO cities (country_id, id, name) VALUES (1, 12, 'hokkaido');
			INSERT INTO cities (country_id, id, name) VALUES (1, 13, 'okinawa');

			INSERT INTO addresses (city_id, id, name) VALUES (11, 111, 'shinjuku');
			INSERT INTO addresses (city_id, id, name) VALUES (11, 112, 'yoyogi');
			INSERT INTO addresses (city_id, id, name) VALUES (11, 113, 'nakano');
			INSERT INTO addresses (city_id, id, name) VALUES (12, 121, 'sapporo');
			INSERT INTO addresses (city_id, id, name) VALUES (12, 122, 'kushiro');
			INSERT INTO addresses (city_id, id, name) VALUES (13, 131, 'naha');
			INSERT INTO addresses (city_id, id, name) VALUES (13, 132, 'miyako');

			INSERT INTO countries (id, name) VALUES (2, 'Somewhere');
			INSERT INTO cities (country_id, id, name) VALUES (2, 21, 'city1');
			INSERT INTO addresses (city_id, id, name) VALUES (21, 211, 'address1');
			INSERT INTO addresses (city_id, id, name) VALUES (21, 212, 'address2');
		`,
		run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
			q := z.Select(
				z.Countries.All(),
				z.Cities.All(),
				z.Addresses.All(),
			).From(
				z.Countries,
			).Joins(
				z.InnerJoin(z.Cities).On(
					z.Cities.CountryID.Eq(z.Countries.ID),
				),
				z.InnerJoin(z.Addresses).On(
					z.Addresses.CityID.Eq(z.Cities.ID),
				),
			).OrderBy(
				z.Countries.ID,
				z.Cities.ID,
				z.Addresses.ID,
			)

			var countries []Country
			var cities map[int][]City
			var addresses map[int][]Address
			err := tx.Query(q).Collect(
				z.Countries.ToUniqSlice(&countries),
				z.Cities.ToUniqSliceMap(&cities).By(z.Countries.ID),
				z.Addresses.ToSliceMap(&addresses).By(z.Cities.ID),
			)
			if err != nil {
				return err
			}

			wantCountries := []Country{
				{1, "Japan", defaultUpdatedAt},
				{2, "Somewhere", defaultUpdatedAt},
			}
			if diff := deep.Equal(countries, wantCountries); diff != nil {
				t.Log(q.Construct())
				return fmt.Errorf("countries diff: %s", diff)
			}

			wantCities := map[int][]City{
				1: []City{
					{11, "tokyo", 1, defaultUpdatedAt},
					{12, "hokkaido", 1, defaultUpdatedAt},
					{13, "okinawa", 1, defaultUpdatedAt},
				},
				2: []City{
					{21, "city1", 2, defaultUpdatedAt},
				},
			}
			if diff := deep.Equal(cities, wantCities); diff != nil {
				t.Log(q.Construct())
				return fmt.Errorf("cities diff: %s", diff)
			}

			wantAddresses := map[int][]Address{
				11: []Address{
					{111, "shinjuku", 11, defaultUpdatedAt},
					{112, "yoyogi", 11, defaultUpdatedAt},
					{113, "nakano", 11, defaultUpdatedAt},
				},
				12: []Address{
					{121, "sapporo", 12, defaultUpdatedAt},
					{122, "kushiro", 12, defaultUpdatedAt},
				},
				13: []Address{
					{131, "naha", 13, defaultUpdatedAt},
					{132, "miyako", 13, defaultUpdatedAt},
				},
				21: []Address{
					{211, "address1", 21, defaultUpdatedAt},
					{212, "address2", 21, defaultUpdatedAt},
				},
			}
			if diff := deep.Equal(addresses, wantAddresses); diff != nil {
				t.Log(q.Construct())
				return fmt.Errorf("addresses diff: %s", diff)
			}

			return nil
		},
	},
	{
		name: "collect records into slice maps through medium table",
		data: `
			INSERT INTO countries (id, name) VALUES (1, 'Japan');

			INSERT INTO cities (country_id, id, name) VALUES (1, 11, 'tokyo');
			INSERT INTO cities (country_id, id, name) VALUES (1, 12, 'hokkaido');
			INSERT INTO cities (country_id, id, name) VALUES (1, 13, 'okinawa');

			INSERT INTO addresses (city_id, id, name) VALUES (11, 111, 'shinjuku');
			INSERT INTO addresses (city_id, id, name) VALUES (11, 112, 'yoyogi');
			INSERT INTO addresses (city_id, id, name) VALUES (11, 113, 'nakano');
			INSERT INTO addresses (city_id, id, name) VALUES (12, 121, 'sapporo');
			INSERT INTO addresses (city_id, id, name) VALUES (12, 122, 'kushiro');
			INSERT INTO addresses (city_id, id, name) VALUES (13, 131, 'naha');
			INSERT INTO addresses (city_id, id, name) VALUES (13, 132, 'miyako');

			INSERT INTO countries (id, name) VALUES (2, 'Somewhere');
			INSERT INTO cities (country_id, id, name) VALUES (2, 21, 'city1');
			INSERT INTO addresses (city_id, id, name) VALUES (21, 211, 'address1');
			INSERT INTO addresses (city_id, id, name) VALUES (21, 212, 'address2');
		`,
		run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
			q := z.Select(
				z.Countries.All(),
				z.Addresses.All(),
			).From(
				z.Countries,
			).Joins(
				z.InnerJoin(z.Cities).On(
					z.Cities.CountryID.Eq(z.Countries.ID),
				),
				z.InnerJoin(z.Addresses).On(
					z.Addresses.CityID.Eq(z.Cities.ID),
				),
			).OrderBy(
				z.Countries.ID,
				z.Cities.ID,
				z.Addresses.ID,
			)

			var countries map[int]Country
			var addresses map[int][]Address
			err := tx.Query(q).Collect(
				z.Countries.ToMap(&countries),
				z.Addresses.ToSliceMap(&addresses).By(z.Countries.ID),
			)
			if err != nil {
				return err
			}

			wantCountries := map[int]Country{
				1: {1, "Japan", defaultUpdatedAt},
				2: {2, "Somewhere", defaultUpdatedAt},
			}
			if diff := deep.Equal(countries, wantCountries); diff != nil {
				t.Log(q.Construct())
				return fmt.Errorf("countries diff: %s", diff)
			}

			wantAddresses := map[int][]Address{
				1: []Address{
					{111, "shinjuku", 11, defaultUpdatedAt},
					{112, "yoyogi", 11, defaultUpdatedAt},
					{113, "nakano", 11, defaultUpdatedAt},
					{121, "sapporo", 12, defaultUpdatedAt},
					{122, "kushiro", 12, defaultUpdatedAt},
					{131, "naha", 13, defaultUpdatedAt},
					{132, "miyako", 13, defaultUpdatedAt},
				},
				2: []Address{
					{211, "address1", 21, defaultUpdatedAt},
					{212, "address2", 21, defaultUpdatedAt},
				},
			}
			if diff := deep.Equal(addresses, wantAddresses); diff != nil {
				t.Log(q.Construct())
				return fmt.Errorf("addresses diff: %s", diff)
			}

			return nil
		},
	},
}
