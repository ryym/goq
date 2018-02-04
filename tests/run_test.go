package tests

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
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

func RunIntegrationTest(t *testing.T, dbName, connStr string) {
	if !ShouldRun(dbName) {
		t.Logf("skip tests for %s", dbName)
		return
	}

	db, err := goq.Open(dbName, connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		t.Fatalf("failed to ping DB: %s", err)
	}

	sql, err := ioutil.ReadFile(filepath.Join("_data", dbName+".sql"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.DB.Exec(string(sql)); err != nil {
		t.Fatal("failed to create tables: %s", err)
	}

	testCases := MakeTestCases(testCtx{dbName})
	RunTestCases(t, db, testCases)
}

func RunTestCases(t *testing.T, db *goq.DB, testCases []testCase) {
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
