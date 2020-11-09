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

	"golang.org/x/crypto/bcrypt"
)

func Initialize() *sql.DB {
	log.Println("Initializing Jablko Database...")

	// Check if database exists
	if _, err := os.Stat("./data/jablko.db"); err == nil {
		log.Println(`Found "jablko.db" in "./data". Using as primary database.`)
	} else if os.IsNotExist(err) {
		log.Println("Database file does not exist. Creating database in \"./data\".")
		createDatabase()
		return nil
	} else {
		log.Println("Issue determining if database file exists. Please check file permisions.")
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

	userTableSQL := `CREATE TABLE users (id INTEGER PRIMARY KEY, username TEXT NOT NULL, password TEXT NOT NULL, firstName TEXT NOT NULL, permissions INTEGER NOT NULL)`
	_, err := newDB.Exec(userTableSQL)
	if err != nil {
		removeDatabase()
		log.Println("FATAL ERROR: Unable to create necessary database. Check file permissions")
		log.Println("Use the following panic information to help with debugging")
		log.Fatal(err.Error())
	}

	// -------------------- Administrative Account --------------------
	log.Println("You must create an administrative account.")
	var username string
	var password string
	var firstName string

	log.Printf("Enter a username:")
	fmt.Scanln(&username)
	log.Printf("Enter a password:")
	fmt.Scanln(&password)
	log.Printf("Enter First Name:")
	fmt.Scanln(&firstName)

	log.Println(username, password)
	adminPassHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		removeDatabase()
		log.Println("FATAL ERROR: Unable to hash administrative password.")
		log.Fatal(err.Error())
	}
	log.Println(string(adminPassHash))

	adminSQL := `INSERT INTO users (username, password, firstName, permissions
) VALUES (?, ?, ?, ?)`
	statement, err := newDB.Prepare(adminSQL)
	if err != nil {
		removeDatabase()
		log.Println("FATAL ERROR: SQL statement for admin account incorrect.")
		log.Fatal(err.Error())
	}

	_, err = statement.Exec(username, adminPassHash, firstName, 2)
	if err != nil {
		removeDatabase()
		log.Println("FATAL ERROR: Unable to insert administrative information into database.")
		log.Fatal(err.Error())
	}
	// -------------------- END Administrative Account --------------------

	// -------------------- Login Sessions --------------------
	log.Println("Creating login sessions table...")

	loginSQL := `CREATE TABLE loginSessions (id INTEGER PRIMARY KEY, cookie TEXT NOT NULL, username TEXT NOT NULL, permissions INTEGER NOT NULL)`
	_, err = newDB.Exec(loginSQL)	
	if err != nil {
		removeDatabase()
		log.Println("FATAL ERROR: Unable to create login sessions table.")
		log.Fatal(err.Error())
	}
	// -------------------- END Login Sessions --------------------
}

func removeDatabase() {
	os.Remove("./data/jablko.db")	
}

func AddUser(db *sql.DB, username string, password string, firstName string, permissions int) error {

	userSQL := `INSERT INTO users (username, password, firstName, permissions) VALUES(?, ?, ?, ?)`

	statement, err := db.Prepare(userSQL)
	if err != nil {
		log.Println("Error in preparing user create SQL statement")
		return err
	}

	_, err = statement.Exec(username, password, firstName, permissions)
	if err != nil {
		log.Println("Error inserting new user in database.")
		return err
	}

	return nil
}
