package Database_Management

import "log"

var RosterDBPath = "../../roster_management"
var DriverDBPath = "../../driver_management"
var DriverAuthDBPath = "../../driver_auth"

var RosterDBInit = "CREATE TABLE IF NOT EXISTS roster (id INTEGER PRIMARY KEY, Drivername TEXT NOT NULL UNIQUE, Rate INTEGER)"
var DriverDBInit = "CREATE TABLE IF NOT EXISTS drivers (id INTEGER PRIMARY KEY, Drivername TEXT NOT NULL UNIQUE, Rate INTEGER)"
var DriverAuthDBInit = "CREATE TABLE IF NOT EXISTS auth (id INTEGER PRIMARY KEY, Username TEXT NOT NULL UNIQUE, Password TEXT)"

func CreateDatabase(db_name string, db_query string) error {

	dbSchema := Database{
		DbName: db_name,
		Query:  db_query,
	}

	err := dbSchema.ExecDB()
	if err != nil {
		log.Panic("Error initialising database: " + db_name)
		return err
	} else {
		log.Println("Database: " + db_name + " created/found")
	}

	return nil

}