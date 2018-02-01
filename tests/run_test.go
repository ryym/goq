package tests

import (
	"context"
	"testing"

	"github.com/ryym/goq"
)

type testCase struct {
	name string
	data string
	run  func(t *testing.T, tx *goq.Tx, z *Builder) error
}

func RunIntegrationTest(t *testing.T, db *goq.DB) {
	builder := NewBuilder(db.Dialect())
	for _, c := range testCases {
		RunTestCase(t, c, builder, db)
	}
}

func RunTestCase(t *testing.T, c testCase, builder *Builder, db *goq.DB) {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = tx.Rollback()
		if err != nil {
			t.Fatal(err)
		}
	}()

	_, err = tx.Tx.Exec(c.data)
	if err != nil {
		t.Fatal(err)
	}

	err = c.run(t, tx, builder)
	if err != nil {
		t.Logf("FAIL: %s", c.name)
		t.Error(err)
	}
}
