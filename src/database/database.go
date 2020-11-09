// database.go: Jablko User Database Management
// Cale Overstreet
// 2020/11/8
// Contains functions that make manipulating the jablko.db 
// easier. Also prevents Jablko Modules from directly 
// modifying the database as the database handle must be
// passed to these functions. This is an effort to prevent 
// malicious module actions.

package database

import (
	"log"
	"os"
	"fmt"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func Initialize() *sql.DB {
	log.Println("Initializing Jablko Database")

	// Check if database exists
	if _, err := os.Stat("./data/jablko.db"); err == nil {
		log.Println(`Found "jablko.db" in "./data". Using as primary database.`)
	} else if os.IsNotExist(err) {
		log.Println("Database file does not exist. Creating database in \"./data\"")
		createDatabase()
		return nil
	} else {
		log.Println("Issue determining if database file exists. Please check file permisions")
	}


	newDB, _ := sql.Open("sqlite3", "./data/jablko.db")		

	log.Println(newDB)
	return newDB
}

func createDatabase() {
	os.Create("./data/jablko.db")

	newDB, _ := sql.Open("sqlite3", "./data/jablko.db")		
	defer newDB.Close()

	log.Println("Creating table \"users\" in database.")

	userTableSQL := `CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT NOT NULL)`
	_, err := newDB.Exec(userTableSQL)
	if err != nil {
		log.Println("FATAL ERROR: Unable to create necessary database. Check file permissions")
		log.Println("Use the following panic information to help with debugging")
		panic(err)
	}

	log.Println("You must create an administrative account.")
	var username string
	log.Printf("Enter a username: ")
	fmt.Scanln(&username)

	var password string
	log.Printf("Enter a password: ")
	fmt.Scanln(&password)

	log.Println(username, password)

	testUserSQL := `INSERT INTO users (name) VALUES ("Cale"), ("Fart"), ("Foo"), ("Bar")`
	_, err = newDB.Exec(testUserSQL)
	if err != nil {
		panic(err)
	}
}
