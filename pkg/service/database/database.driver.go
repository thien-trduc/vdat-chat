package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() *sql.DB {
	conn := "postgres://postgres:postgres@localhost:5432/dchat"

	connectionStr := os.Getenv("DATABASE_URL")
	if len(connectionStr) > 0 {
		conn = connectionStr
	}
	conn = conn + "?sslmode=disable"
	log.Println("Connect string", conn)
	db, err := sql.Open("postgres", conn)

	if err != nil {
		fmt.Printf("Fail to openDB: %v \n", err)
	}
	DB = db
	return DB
}
