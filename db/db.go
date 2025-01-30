package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var Db *sql.DB

func InitDatabase() error {

	// databse config
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	// MySQL connection string
	connStr := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

	// Open the Postgres connection
	var err error
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	// defer db.Close()

	// Ensure the connection is working
	err = Db.Ping()
	if err != nil {
		return err
	}

	// database initialized with no error
	return nil
}
