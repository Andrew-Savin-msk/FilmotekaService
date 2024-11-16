package pgstore_test

import (
	"os"
	"testing"
)

var (
	dbPath string
)

func TestMain(m *testing.M) {
	if path := os.Getenv("DB_PATH_TEST"); path == "" {
		dbPath = "postgres://postgres:Sassassa12@localhost/restapi_test?sslmode=disable"
	} else {
		dbPath = path
	}
}
