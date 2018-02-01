package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/ryym/goq"
)

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
}
