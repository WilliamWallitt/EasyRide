package Database_Management

import "log"

func CreateDatabases() {

	databaseMapping := map[string]string {
		"./roster_management": "CREATE TABLE IF NOT EXISTS roster (id INTEGER PRIMARY KEY, Drivername TEXT, Rate INTEGER )",
		"./driver_management": "CREATE TABLE IF NOT EXISTS drivers (id INTEGER PRIMARY KEY, Drivername TEXT, Rate INTEGER )",
		"./driver_auth": "CREATE TABLE IF NOT EXISTS auth (id INTEGER PRIMARY KEY, Username TEXT, Password TEXT)",
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
		} else {
			log.Println("Database: " + key + " created")
		}
	}

}