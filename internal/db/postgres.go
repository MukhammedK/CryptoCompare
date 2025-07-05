package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func ConnectPostgres() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=postgres password=1234 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Connected to PostgreSQL!")
	return db, nil
}
