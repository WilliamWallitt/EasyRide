package Database_Management

import "log"

var RosterDBPath = "./Roster_Management_/roster_management"
var DriverDBPath = "./Driver_Management_/driver_management"
var DriverAuthDBPath = "./Driver_Authentication/driver_auth"

func CreateDatabases() error {

	databaseMapping := map[string]string {
		RosterDBPath: "CREATE TABLE IF NOT EXISTS roster (id INTEGER PRIMARY KEY, Drivername TEXT, Rate INTEGER)",
		DriverDBPath: "CREATE TABLE IF NOT EXISTS drivers (id INTEGER PRIMARY KEY, Drivername TEXT, Rate INTEGER)",
		DriverAuthDBPath: "CREATE TABLE IF NOT EXISTS auth (id INTEGER PRIMARY KEY, Username TEXT NOT NULL UNIQUE, Password TEXT)",
	}

	dbSchema := Database{
		DbName: "",
		Query:  "",
	}

	for key, value := range databaseMapping {
		dbSchema.DbName, dbSchema.Query = key, value
		err := dbSchema.ExecDB()
		if err != nil {
			log.Panic("Error initialising database: " + key)
			return err
		} else {
			log.Println("Database: " + key + " created/found")
		}
	}

	return nil

}