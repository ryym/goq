package tests

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/ryym/goq"
	"github.com/ryym/goq/cllct"
)

var defaultUpdatedAt = time.Date(2000, time.January, 1, 9, 0, 0, 0, time.UTC)

var testCases = []testCase{
	{
		name: "check all collectors are not broken",
		data: `
			INSERT INTO countries (id, name) VALUES (1, 'Japan');
			INSERT INTO countries (id, name) VALUES (2, 'US');
			INSERT INTO cities (country_id, id, name) VALUES (1, 1, 'chiba');
			INSERT INTO cities (country_id, id, name) VALUES (1, 2, 'gunma');
			INSERT INTO cities (country_id, id, name) VALUES (2, 3, 'miami');
			INSERT INTO cities (country_id, id, name) VALUES (2, 4, 'reno');
		`,
		run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
			q := z.Select(
				z.Countries.ID.As("country_id"),
				z.Cities.ID.As("city_id"),
			).From(z.Countries, z.Cities).Where(
				z.Countries.ID.Eq(z.Cities.CountryID),
			).OrderBy(z.Countries.ID, z.Cities.ID)

			type myCityT struct{ CityID int }
			var (
				myCity       myCityT
				row          map[string]interface{}
				rows         []map[string]interface{}
				myCityMap    map[int]myCityT
				myCities     []myCityT
				myCitiesMap  map[int][]myCityT
				city         City
				cities       []City
				cityMap      map[int]City
				citiesMap    map[int][]City
				countriesMap map[int][]Country
			)
			var keyStore int

			cases := []struct {
				cllct cllct.Collector
				got   interface{}
				want  interface{}
			}{
				{
					cllct: z.ToElem(&myCity),
					got:   &myCity,
					want:  myCityT{1},
				},
				{
					cllct: z.ToRowMap(&row),
					got:   &row,
					want: map[string]interface{}{
						"country_id": int64(1),
						"city_id":    int64(1),
					},
				},
				{
					cllct: z.ToRowMapSlice(&rows),
					got:   &rows,
					want: []map[string]interface{}{
						{"country_id": int64(1), "city_id": int64(1)},
						{"country_id": int64(1), "city_id": int64(2)},
						{"country_id": int64(2), "city_id": int64(3)},
						{"country_id": int64(2), "city_id": int64(4)},
					},
				},
				{
					cllct: z.ToMap(&myCityMap).By(z.Cities.ID),
					got:   &myCityMap,
					want: map[int]myCityT{
						1: {1}, 2: {2},
						3: {3}, 4: {4},
					},
				},
				{
					cllct: z.ToSlice(&myCities),
					got:   &myCities,
					want:  []myCityT{{1}, {2}, {3}, {4}},
				},
				{
					cllct: z.ToSliceMap(&myCitiesMap).ByWith(&keyStore, z.Countries.ID),
					got:   &myCitiesMap,
					want: map[int][]myCityT{
						1: []myCityT{{1}, {2}},
						2: []myCityT{{3}, {4}},
					},
				},
				{
					cllct: z.Cities.ToElem(&city),
					got:   &city,
					want:  City{ID: 1},
				},
				{
					cllct: z.Cities.ToSlice(&cities),
					got:   &cities,
					want: []City{
						{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4},
					},
				},
				{
					cllct: z.Cities.ToMap(&cityMap),
					got:   &cityMap,
					want: map[int]City{
						1: {ID: 1}, 2: {ID: 2},
						3: {ID: 3}, 4: {ID: 4},
					},
				},
				{
					cllct: z.Cities.ToSliceMap(&citiesMap).ByWith(&keyStore, z.Countries.ID),
					got:   &citiesMap,
					want: map[int][]City{
						1: []City{{ID: 1}, {ID: 2}},
						2: []City{{ID: 3}, {ID: 4}},
					},
				},
				{
					cllct: z.Countries.ToUniqSliceMap(&countriesMap).By(z.Countries.ID),
					got:   &countriesMap,
					want: map[int][]Country{
						1: []Country{{ID: 1}},
						2: []Country{{ID: 2}},
					},
				},
			}

			var succeed = true
			for _, c := range cases {
				var err error
				switch cl := c.cllct.(type) {
				case cllct.ListCollector:
					err = tx.Query(q).Collect(cl)
				case cllct.SingleCollector:
					err = tx.Query(q).First(cl)
				}
				if err != nil {
					t.Errorf("%s: %s", reflect.TypeOf(c.cllct), err)
					succeed = false
					continue
				}
				got := reflect.ValueOf(c.got).Elem().Interface()
				if diff := deep.Equal(got, c.want); diff != nil {
					t.Errorf("%s: %s", reflect.TypeOf(c.cllct), diff)
					succeed = false
				}
			}

			if succeed {
				return nil
			} else {
				return errors.New("some collector does not work")
			}
		},
	},
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
	{
		name: "collect records into non-model struct slices",
		data: `
			INSERT INTO cities (id, name, country_id) VALUES (10, 'newyork', 5);
			INSERT INTO cities (id, name, country_id) VALUES (12, 'chicago', 5);
			INSERT INTO cities (id, name, country_id) VALUES (14, 'seattle', 5);
		`,
		run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
			type myCity struct {
				Code      int
				CityName  string
				CountryID int
			}

			q := z.Select(
				z.Cities.ID.As("code"),
				z.Cities.Name.As("city_name"),
				z.Cities.CountryID,
			).From(z.Cities).OrderBy(z.Cities.ID)

			var cities []myCity
			err := tx.Query(q).Collect(z.ToSlice(&cities))
			if err != nil {
				return err
			}

			want := []myCity{
				{Code: 10, CityName: "newyork", CountryID: 5},
				{Code: 12, CityName: "chicago", CountryID: 5},
				{Code: 14, CityName: "seattle", CountryID: 5},
			}
			if diff := deep.Equal(cities, want); diff != nil {
				t.Log(q.Construct())
				return fmt.Errorf("%s", diff)
			}
			return nil
		},
	},
	{
		name: "collect records into map slices",
		data: `
			INSERT INTO cities (id, name, country_id) VALUES (10, 'newyork', 5);
			INSERT INTO cities (id, name, country_id) VALUES (12, 'chicago', 5);
			INSERT INTO cities (id, name, country_id) VALUES (14, 'seattle', 5);
		`,
		run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
			q := z.Select(z.Cities.All()).From(z.Cities)

			var cities []map[string]interface{}
			err := tx.Query(q).Collect(z.ToRowMapSlice(&cities))
			if err != nil {
				return err
			}

			want := []map[string]interface{}{
				{"id": int64(10), "name": "newyork", "country_id": int64(5), "updated_at": defaultUpdatedAt},
				{"id": int64(12), "name": "chicago", "country_id": int64(5), "updated_at": defaultUpdatedAt},
				{"id": int64(14), "name": "seattle", "country_id": int64(5), "updated_at": defaultUpdatedAt},
			}
			if diff := deep.Equal(cities, want); diff != nil {
				t.Log(cities)
				return fmt.Errorf("%s", diff)
			}
			return nil
		},
	},
	{
		name: "collect into models and maps",
		data: `
			INSERT INTO countries (id, name) values (1, 'Japan');
			INSERT INTO countries (id, name) values (2, 'Brazil');
			INSERT INTO countries (id, name) values (3, 'Nowhere');

			INSERT INTO cities (country_id, name) values (1, 'japan-city');
			INSERT INTO cities (country_id, name) values (1, 'japan-city');
			INSERT INTO cities (country_id, name) values (2, 'brazil-city');
			INSERT INTO cities (country_id, name) values (2, 'brazil-city');
			INSERT INTO cities (country_id, name) values (2, 'brazil-city');
			INSERT INTO cities (country_id, name) values (2, 'brazil-city');
		`,
		run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
			q := z.Select(
				z.Countries.ID,
				z.Countries.Name,
				z.Count(z.Cities.ID).As("cities_count"),
				z.VarT(10, "int").Add(20).As("population"),
			).From(z.Countries).Joins(
				z.LeftJoin(z.Cities).On(
					z.Cities.CountryID.Eq(z.Countries.ID),
				),
			).GroupBy(
				z.Countries.ID, z.Countries.Name,
			).OrderBy(z.Countries.ID)

			var countries []Country
			var extras []map[string]interface{}
			err := tx.Query(q).Collect(
				z.Countries.ToSlice(&countries),
				z.ToRowMapSlice(&extras),
			)
			if err != nil {
				t.Log(q.Construct())
				return err
			}

			wantCountries := []Country{
				{ID: 1, Name: "Japan"},
				{ID: 2, Name: "Brazil"},
				{ID: 3, Name: "Nowhere"},
			}
			if diff := deep.Equal(countries, wantCountries); diff != nil {
				return fmt.Errorf("countries diff: %s", diff)
			}

			wantExtras := []map[string]interface{}{
				{"cities_count": int64(2), "population": int64(30)},
				{"cities_count": int64(4), "population": int64(30)},
				{"cities_count": int64(0), "population": int64(30)},
			}
			if diff := deep.Equal(wantExtras, extras); diff != nil {
				return fmt.Errorf("extras diff: %s", diff)
			}

			return nil
		},
	},
	{
		name: "collect records using custom name helper",
		data: `
			INSERT INTO technologies VALUES (1, 'AI', 'what!?');
			INSERT INTO technologies VALUES (2, 'self-driving car', 'cool!');
		`,
		run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
			q := z.Select(z.Techs.All()).From(z.Techs).OrderBy(z.Techs.Desc)

			var techs []Tech
			err := tx.Query(q).Collect(z.Techs.ToSlice(&techs))
			if err != nil {
				return err
			}

			wantTechs := []Tech{
				{2, "self-driving car", "cool!"},
				{1, "AI", "what!?"},
			}
			if diff := deep.Equal(techs, wantTechs); diff != nil {
				return fmt.Errorf("techs diff: %s", diff)
			}
			return nil
		},
	},
	{
		name: "use self-join by specifying different aliases",
		data: `
			INSERT INTO technologies VALUES (1, 'VR', '');
			INSERT INTO technologies VALUES (2, 'AR', '');
			INSERT INTO technologies VALUES (3, 'MR', '');
		`,
		run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
			t1 := z.Techs.As("t1")
			t2 := z.Techs.As("t2")
			q := z.Select(
				t1.Name.As("name1"), t2.Name.As("name2"),
			).From(t1, t2).Where(
				t1.Name.Neq("VR"),
			).OrderBy(t1.Name, t2.Name)

			var techs []map[string]interface{}
			err := tx.Query(q).Collect(z.ToRowMapSlice(&techs))
			if err != nil {
				return err
			}

			wantTechs := []map[string]interface{}{
				{"name1": "AR", "name2": "AR"},
				{"name1": "AR", "name2": "MR"},
				{"name1": "AR", "name2": "VR"},
				{"name1": "MR", "name2": "AR"},
				{"name1": "MR", "name2": "MR"},
				{"name1": "MR", "name2": "VR"},
			}
			if diff := deep.Equal(techs, wantTechs); diff != nil {
				return fmt.Errorf("techs diff: %s", diff)
			}
			return nil
		},
	},
}
