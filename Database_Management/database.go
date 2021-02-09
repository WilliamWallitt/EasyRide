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

//query: "CREATE TABLE IF NOT EXISTS roster (id INTEGER PRIMARY KEY, Username TEXT, Age INTEGER )",
//query: "INSERT INTO roster (Username, Age) VALUES ('William', 21)"
//	"SELECT id, Username, Age FROM roster WHERE id=(1)"


//func main() {

//
//	setup := Database {
//		dbName: "./Database_Management/test.db",
//		query: "SELECT id, Username, Age FROM roster WHERE id=(1)",
//	}
//	//[]string{}
//	rows, err := setup.QueryDB()
//	if err != nil {
//		log.Println(err)
//	} else {
//		var id int
//		var user string
//		var age int
//		for rows.Next() {
//			err := rows.Scan(&id, &user, &age)
//			if err != nil {
//				log.Println(err)
//			}
//			fmt.Println(id, user, age)
//		}
//		log.Println("Added")
//
//	}
//
//}
