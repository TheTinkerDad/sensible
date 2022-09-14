package settings

// TODO: THIS FILE IS JUST A SKELETON RIGHT NOW! Known issues:
// "Save" has hardcoded values
// "Load" keep overwriting the domain section with data[domain] = make(map[string]string)

import (
	"database/sql"
	"log"

	// Needed to work with SQLite DB files
	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB = nil

// domain:key:value
var data map[string]map[string]string = nil

func init() {

	log.Println("Opening settings database...")

	db, err := sql.Open("sqlite3", "./sensible.db")
	if err != nil {
		log.Fatal(err)
	}

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS configuration (id INTEGER PRIMARY KEY, domain TEXT, key TEXT, value TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	database = db

	Load()
}

// Save Saves the current settings
func Save() {

	statement, err := database.Prepare("INSERT INTO configuration (domain, key, value) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec("mqtt", "broker-hostname", "192.168.1.4")
}

// Load Loads the current settings
func Load() {

	data = make(map[string]map[string]string)

	rows, err := database.Query("SELECT domain, key, value FROM configuration")
	if err != nil {
		log.Fatal(err)
	}

	var domain string
	var key string
	var value string
	for rows.Next() {
		rows.Scan(&domain, &key, &value)
		log.Println(domain + "." + key + ": " + value)
		data[domain] = make(map[string]string)
		data[domain][key] = value
	}
}

// Get Returns a configuration parameter
func Get(domain string, key string) string {

	return data[domain][key]
}

// EnsureOk Checks if the loaded configuration is intact
func EnsureOk() {

}
