package tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ryym/goq"
)

func TestPostgres(t *testing.T) {
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5433"
	}
	conn := fmt.Sprintf("port=%s user=goq sslmode=disable", port)
	db, err := goq.Open("postgres", conn)

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}

	sql, err := ioutil.ReadFile(filepath.Join("_data", "postgres.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.DB.Exec(string(sql)); err != nil {
		t.Fatal(err)
	}

	RunIntegrationTest(t, db)
}
