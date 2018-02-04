package tests

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestPostgres(t *testing.T) {
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5433"
	}
	conn := fmt.Sprintf("port=%s user=goq sslmode=disable", port)
	RunIntegrationTest(t, "postgres", conn)
}
