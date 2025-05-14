package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

const DATABASE_NAME = "./procurement.db" // Relative to backend directory, where db_init.go creates it

var DB *sql.DB // Global DB connection pool

// InitDB initializes the database connection
func InitDB() error {
	var err error
	// The database file is expected to be in the same directory as the executable
	// or where db_init.go creates it, which is the 'backend' directory.
	db, err := sql.Open("sqlite3", DATABASE_NAME+"?_foreign_keys=on")
	if err != nil {
		log.Printf("Error opening database specified by DATABASE_NAME ('%s'): %v", DATABASE_NAME, err)
		return err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		log.Printf("Error pinging database: %v", err)
		return err
	}
	DB = db
	log.Println("Successfully connected to SQLite database from db.go!")
	return nil
}

// GetDB returns the current database connection pool
func GetDB() *sql.DB {
	if DB == nil {
		log.Println("Warning: GetDB called before InitDB or InitDB failed.")
	}
	return DB
}
