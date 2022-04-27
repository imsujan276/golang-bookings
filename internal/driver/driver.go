package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB holds the databse connection pool
type DB struct {
	SQL *sql.DB
}

var dbCon = &DB{}

const (
	maxOpenDbCon  = 10
	maxIdleDbCon  = 5
	maxDbLifetime = 5 * time.Minute
)

// ConnectSQL creates database pool for postgres
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		// exits the program
		panic(err)
	}
	d.SetMaxOpenConns(maxOpenDbCon)
	d.SetMaxIdleConns(maxIdleDbCon)
	d.SetConnMaxLifetime(maxDbLifetime)

	dbCon.SQL = d

	if err = testDB(d); err != nil {
		return nil, err
	}
	return dbCon, nil

}

// NewDatabase creates new database for the application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// testDB tries to ping Database
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}
