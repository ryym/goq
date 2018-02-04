package tests

import (
	"context"
	"os"
	"testing"

	"github.com/ryym/goq"
)

type testCase struct {
	name string
	data string
	run  func(t *testing.T, tx *goq.Tx, z *Builder) error
	only bool
}

func ShouldRun(dbName string) bool {
	target := os.Getenv("DB")
	return target == "" || target == dbName
}

func RunIntegrationTest(t *testing.T, db *goq.DB) {
	var targets []testCase
	for _, c := range testCases {
		if c.only {
			targets = append(targets, c)
		}
	}
	if len(targets) == 0 {
		targets = testCases
	}

	builder := NewBuilder(db.Dialect())
	for _, c := range targets {
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
