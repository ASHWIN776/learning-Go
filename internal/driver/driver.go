package driver

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB holds the Database connection pool as one of its variables
type DB struct {
	SQL *sql.DB
}

/*
	* Reason to use a struct?
	In the future maybe I want to use a difference db, so I can add that into this struct
*/

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

// Returns the instance of DB type which has the connection pool with the above constants attached to it
func ConnectSQL(dsn string) (*DB, error) {
	// Get *sql.DB from New Database function
	db, err := NewDatabase(dsn)

	if err != nil {
		log.Println("could not connect to DB")
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDbLifetime)

	// Test the database again
	err = TestDatabase(db)
	if err != nil {
		log.Println("could not connect to db after setting constants")
		return nil, err
	}

	dbConn.SQL = db

	return dbConn, nil
}

// Tries to ping the database and returns whatever gets back(type of error)
func TestDatabase(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}

// Creates the connection pool and an error if any
func NewDatabase(dsn string) (*sql.DB, error) {
	// Get connection pool
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	// Test Connection pool
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil

}
