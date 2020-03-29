package tests

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ryym/goq"
)

func MakeTestCases(ctx testCtx) []testCase {
	return []testCase{
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
					cllct goq.Collector
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
							"country_id": ctx.rawInt(1),
							"city_id":    ctx.rawInt(1),
						},
					},
					{
						cllct: z.ToRowMapSlice(&rows),
						got:   &rows,
						want: []map[string]interface{}{
							{"country_id": ctx.rawInt(1), "city_id": ctx.rawInt(1)},
							{"country_id": ctx.rawInt(1), "city_id": ctx.rawInt(2)},
							{"country_id": ctx.rawInt(2), "city_id": ctx.rawInt(3)},
							{"country_id": ctx.rawInt(2), "city_id": ctx.rawInt(4)},
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
					case goq.ListCollector:
						err = tx.Query(q).Collect(cl)
					case goq.SingleCollector:
						err = tx.Query(q).First(cl)
					}
					if err != nil {
						t.Errorf("%s: %s", reflect.TypeOf(c.cllct), err)
						succeed = false
						continue
					}
					got := reflect.ValueOf(c.got).Elem().Interface()
					if diff := cmp.Diff(got, c.want); diff != "" {
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
				if diff := cmp.Diff(countries, wantCountries); diff != "" {
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
				if diff := cmp.Diff(cities, wantCities); diff != "" {
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
				if diff := cmp.Diff(countries, wantCountries); diff != "" {
					t.Log(q.Construct())
					return fmt.Errorf("countries diff: %s", diff)
				}

				wantCities := []City{
					{Name: "hokkaido", CountryID: 8},
					{Name: "okinawa", CountryID: 8},
					{Name: "tokyo", CountryID: 8},
					{Name: "foo", CountryID: 20},
				}
				if diff := cmp.Diff(cities, wantCities); diff != "" {
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
					{1, "Japan"},
					{2, "Somewhere"},
				}
				if diff := cmp.Diff(countries, wantCountries); diff != "" {
					return fmt.Errorf("countries diff: %s", diff)
				}

				wantCities := map[int][]City{
					1: []City{
						{11, "tokyo", 1},
						{12, "hokkaido", 1},
						{13, "okinawa", 1},
					},
					2: []City{
						{21, "city1", 2},
					},
				}
				if diff := cmp.Diff(cities, wantCities); diff != "" {
					return fmt.Errorf("cities diff: %s", diff)
				}

				wantAddresses := map[int][]Address{
					11: []Address{
						{111, "shinjuku", 11},
						{112, "yoyogi", 11},
						{113, "nakano", 11},
					},
					12: []Address{
						{121, "sapporo", 12},
						{122, "kushiro", 12},
					},
					13: []Address{
						{131, "naha", 13},
						{132, "miyako", 13},
					},
					21: []Address{
						{211, "address1", 21},
						{212, "address2", 21},
					},
				}
				if diff := cmp.Diff(addresses, wantAddresses); diff != "" {
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
					1: {1, "Japan"},
					2: {2, "Somewhere"},
				}
				if diff := cmp.Diff(countries, wantCountries); diff != "" {
					t.Log(q.Construct())
					return fmt.Errorf("countries diff: %s", diff)
				}

				wantAddresses := map[int][]Address{
					1: []Address{
						{111, "shinjuku", 11},
						{112, "yoyogi", 11},
						{113, "nakano", 11},
						{121, "sapporo", 12},
						{122, "kushiro", 12},
						{131, "naha", 13},
						{132, "miyako", 13},
					},
					2: []Address{
						{211, "address1", 21},
						{212, "address2", 21},
					},
				}
				if diff := cmp.Diff(addresses, wantAddresses); diff != "" {
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
				if diff := cmp.Diff(cities, want); diff != "" {
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
					{
						"id":         ctx.rawInt(10),
						"name":       ctx.rawStr("newyork"),
						"country_id": ctx.rawInt(5),
					},
					{
						"id":         ctx.rawInt(12),
						"name":       ctx.rawStr("chicago"),
						"country_id": ctx.rawInt(5),
					},
					{
						"id":         ctx.rawInt(14),
						"name":       ctx.rawStr("seattle"),
						"country_id": ctx.rawInt(5),
					},
				}
				if diff := cmp.Diff(cities, want); diff != "" {
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
				if diff := cmp.Diff(countries, wantCountries); diff != "" {
					return fmt.Errorf("countries diff: %s", diff)
				}

				wantExtras := []map[string]interface{}{
					// These values are always `int64` so we don't need to use `ctx.rawInt`.
					// RDB seems to be able to infer their types from `COUNT` function
					// and literal values.
					{"cities_count": int64(2), "population": int64(30)},
					{"cities_count": int64(4), "population": int64(30)},
					{"cities_count": int64(0), "population": int64(30)},
				}
				if diff := cmp.Diff(wantExtras, extras); diff != "" {
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
				if diff := cmp.Diff(techs, wantTechs); diff != "" {
					return fmt.Errorf("techs diff: %s", diff)
				}
				return nil
			},
		},
		{
			name: "return ErrNoRows for empty result when First is used",
			data: ``,
			run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
				q := z.Select(z.Cities.All()).From(z.Cities)
				var city City
				err := tx.Query(q).First(z.Cities.ToElem(&city))
				if err == nil {
					return errors.New("no error returned")
				}
				if !errors.Is(err, goq.ErrNoRows) {
					return fmt.Errorf("unknown error returned: %v", err)
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

				type techPair struct {
					Name1 string
					Name2 string
				}

				var techs []techPair
				err := tx.Query(q).Collect(z.ToSlice(&techs))
				if err != nil {
					return err
				}

				wantTechs := []techPair{
					{"AR", "AR"},
					{"AR", "MR"},
					{"AR", "VR"},
					{"MR", "AR"},
					{"MR", "MR"},
					{"MR", "VR"},
				}
				if diff := cmp.Diff(techs, wantTechs); diff != "" {
					return fmt.Errorf("techs diff: %s", diff)
				}
				return nil
			},
		},
		{
			name: "ToMap for duplicated records selects last one",
			data: `
			INSERT INTO cities (country_id, id, name) VALUES (1, 1, 'foo');
			INSERT INTO cities (country_id, id, name) VALUES (1, 2, 'foo');
			`,
			run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
				q := z.Select(z.Cities.ID, z.Cities.Name).From(z.Cities).OrderBy(z.Cities.ID)

				var cities map[string]City
				err := tx.Query(q).Collect(z.ToMap(&cities).By(z.Cities.Name))
				if err != nil {
					return err
				}

				want := map[string]City{"foo": {ID: 2, Name: "foo"}}
				if diff := cmp.Diff(cities, want); diff != "" {
					t.Log(q.Construct())
					return fmt.Errorf("%s", diff)
				}
				return nil
			},
		},
		{
			name: "Insert new records",
			data: "",
			run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
				city := City{
					ID:        80,
					Name:      "city-80",
					CountryID: 33,
				}

				q := z.InsertInto(z.Cities).Values(city)
				_, err := tx.Exec(q)
				if err != nil {
					return err
				}

				var got City
				tx.Query(z.Select(z.Cities.All()).From(z.Cities)).First(
					z.Cities.ToElem(&got),
				)
				if diff := cmp.Diff(city, got); diff != "" {
					return fmt.Errorf("%s", diff)
				}
				return nil
			},
		},
		{
			name: "Insert new records by specifying columns",
			data: "",
			run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
				tech := Tech{
					ID:   1,
					Name: "Go",
				}

				q := z.InsertInto(z.Techs, z.Techs.ID, z.Techs.Name).Values(tech)
				_, err := tx.Exec(q)
				if err != nil {
					return err
				}

				var got Tech
				tx.Query(z.Select(z.Techs.All()).From(z.Techs)).First(
					z.Techs.ToElem(&got),
				)
				if diff := cmp.Diff(tech, got); diff != "" {
					return fmt.Errorf("%s", diff)
				}
				return nil
			},
		},
		{
			name: "Update specific records",
			data: `
			INSERT INTO cities (country_id, id, name) VALUES (1, 1, 'a');
			INSERT INTO cities (country_id, id, name) VALUES (1, 2, 'b');
			INSERT INTO cities (country_id, id, name) VALUES (1, 3, 'c');
			`,
			run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
				_, err := tx.Exec(
					z.Update(z.Cities).Set(goq.Values{
						z.Cities.CountryID: 50,
						z.Cities.Name:      "x",
					}).Where(z.Cities.ID.Lt(3)),
				)
				if err != nil {
					return err
				}

				want := []City{
					{ID: 1, CountryID: 50, Name: "x"},
					{ID: 2, CountryID: 50, Name: "x"},
					{ID: 3, CountryID: 1, Name: "c"},
				}

				var got []City
				tx.Query(
					z.Select(z.Cities.All()).From(z.Cities).OrderBy(z.Cities.ID),
				).Collect(z.Cities.ToSlice(&got))

				if diff := cmp.Diff(want, got); diff != "" {
					return fmt.Errorf("%s", diff)
				}
				return nil
			},
		},
		{
			name: "Update one record using struct",
			data: `
			INSERT INTO cities (country_id, id, name) VALUES (1, 1, 'a');
			INSERT INTO cities (country_id, id, name) VALUES (1, 2, 'b');
			`,
			run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
				_, err := tx.Exec(
					z.Update(z.Cities).Elem(City{
						ID:        1,
						Name:      "x",
						CountryID: 8,
					}),
				)
				if err != nil {
					return err
				}

				want := []City{
					{ID: 1, CountryID: 8, Name: "x"},
					{ID: 2, CountryID: 1, Name: "b"},
				}

				var got []City
				tx.Query(
					z.Select(z.Cities.All()).From(z.Cities).OrderBy(z.Cities.ID),
				).Collect(z.Cities.ToSlice(&got))

				if diff := cmp.Diff(want, got); diff != "" {
					return fmt.Errorf("%s", diff)
				}
				return nil
			},
		},
		{
			name: "Delete records",
			data: `
			INSERT INTO cities (country_id, id, name) VALUES (1, 1, 'a');
			INSERT INTO cities (country_id, id, name) VALUES (2, 2, 'b');
			INSERT INTO cities (country_id, id, name) VALUES (2, 3, 'c');
			`,
			run: func(t *testing.T, tx *goq.Tx, z *Builder) error {
				_, err := tx.Exec(
					z.DeleteFrom(z.Cities).Where(z.Cities.CountryID.Eq(2)),
				)
				if err != nil {
					return err
				}

				want := []City{
					{ID: 1, CountryID: 1, Name: "a"},
				}

				var got []City
				tx.Query(
					z.Select(z.Cities.All()).From(z.Cities).OrderBy(z.Cities.ID),
				).Collect(z.Cities.ToSlice(&got))

				if diff := cmp.Diff(want, got); diff != "" {
					return fmt.Errorf("%s", diff)
				}
				return nil
			},
		},
	}
}
