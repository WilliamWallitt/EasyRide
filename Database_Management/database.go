package Database_Management

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DbName, Query string
}


// delete, update, create db methods

func (db *Database) ExecDB() error {
	_db, err := sql.Open("sqlite3", db.DbName)
	if err != nil {
		return err
	}

	statement, err := _db.Prepare(db.Query)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	err = statement.Close()
	if err != nil {
		return err
	}
	return nil
}

// get / get all methods

func (db *Database) QueryDB() (*sql.Rows, error) {
	_db, err := sql.Open("sqlite3", db.DbName)
	if err != nil {
		return nil, err
	}
	statement, err := _db.Prepare(db.Query)
	if err != nil {
		return nil, err
	}

	rows, err := statement.Query()
	if err != nil {
		return nil, err
	}
	err = statement.Close()
	if err != nil {
		return nil, err
	}
	return rows, nil
}
