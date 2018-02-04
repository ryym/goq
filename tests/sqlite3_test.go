package tests

import (
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// XXX: Currently you need to install SQLite3 to run this test.
// https://sqlite.org/index.html
func TestSQLite3(t *testing.T) {
	RunIntegrationTest(t, DB_SQLITE3, "./_sqlite.db")
	defer os.Remove("./_sqlite.db")
}
