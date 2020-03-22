package tests

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func TestPostgres(t *testing.T) {
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5434"
	}
	conn := fmt.Sprintf("port=%s user=goq sslmode=disable", port)
	RunIntegrationTest(t, DB_POSTGRES, conn)
}

func TestMySQL(t *testing.T) {
	port := os.Getenv("MYSQL_PORT")
	if port == "" {
		port = "3307"
	}
	conn := fmt.Sprintf("root:root@(:%s)/goq?multiStatements=true", port)
	RunIntegrationTest(t, DB_MYSQL, conn)
}
