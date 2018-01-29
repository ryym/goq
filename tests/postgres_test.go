package tests

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestPostgres(t *testing.T) {
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}
	conn := fmt.Sprintf("port=%s user=goq sslmode=disable", port)
	db, err := sql.Open("postgres", conn)

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}

	if _, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id integer)"); err != nil {
		t.Fatal(err)
	}

	rows, err := db.Query("SELECT COUNT(*) FROM users")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	rows.Next()

	var count int
	rows.Scan(&count)
	if count != 0 {
		t.Errorf("[users table count] want: 0, got: %d", count)
	}
}
