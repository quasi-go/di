package sql

import "fmt"

// This file create a fake stub of the database/sql package

type DB struct {
	// this would be represented a connection to a database, something like sql.Db
}

func (db *DB) Query(query string, args ...any) {
	fmt.Println("RECEIVED QUERY: ", query, args)
}

func (db *DB) Ping() error {
	fmt.Println("PINGED CONNECTION")
	return nil
}

func Open(driver string, dsn string) (*DB, error) {
	fmt.Println("OPENING CONNECTION WITH DRIVER ", driver, "TO", dsn)
	return &DB{}, nil
}
