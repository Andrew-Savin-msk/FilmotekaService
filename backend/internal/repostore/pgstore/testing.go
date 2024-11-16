package pgstore

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

func TestActorMap() map[string]interface{} {

	return map[string]interface{}{
		"name":      "John",
		"gender":    "man",
		"birthdate": time.Now(),
	}
}

func TestStore(dbURL string) (*sql.DB, func(...string)) {

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db, func(s ...string) {
		if len(s) > 0 {
			db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", strings.Join(s, ", ")))
		}
	}
}
