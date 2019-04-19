package db

import (
    "os"
	"database/sql"
	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
    user := os.Getenv("POSTGRES_USER")
    dbname := os.Getenv("POSTGRES_DB")
    passwd := os.Getenv("POSTGRES_PASSWORD")

	connStr := "host=db user="+user+" dbname="+dbname+" password="+passwd+" sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}
