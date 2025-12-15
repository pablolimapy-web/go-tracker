package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewConnection(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("db not reachable:", err)
	}

	return db
}
