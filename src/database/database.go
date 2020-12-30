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
	"os"
	"fmt"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"golang.org/x/crypto/bcrypt"

	"github.com/ccoverstreet/Jablko/types"
	"github.com/ccoverstreet/Jablko/src/jlog"
)

type JablkoDB struct {
	Db *sql.DB
}

func Initialize() *JablkoDB {
	newDB := new(JablkoDB)
	jlog.Println("Initializing Jablko Database...")

	// Check if database exists
	if _, err := os.Stat("./data/jablko.db"); err == nil {
		jlog.Println(`Found "jablko.db" in "./data". Using as primary database.`)
	} else if os.IsNotExist(err) {
		jlog.Warnf("Database file does not exist. Creating database in \"./data\".\n")
		createDatabase()
	} else {
		jlog.Errorf("Issue determining if database file exists. Please check file permisions.\n")
	}

	newDB.Db, _ = sql.Open("sqlite3", "./data/jablko.db")		

	return newDB
}

func createDatabase() {
	os.Create("./data/jablko.db")

	newDB, _ := sql.Open("sqlite3", "./data/jablko.db")		
	defer newDB.Close()

	jlog.Println("Creating table \"users\" in database.")

	userTableSQL := `CREATE TABLE users (id INTEGER PRIMARY KEY, username TEXT NOT NULL, password TEXT NOT NULL, firstName TEXT NOT NULL, permissions INTEGER NOT NULL)`
	_, err := newDB.Exec(userTableSQL)
	if err != nil {
		removeDatabase()
		jlog.Errorf("Unable to create necessary database. Check file permissions.\n")
		jlog.Println(err)
		panic(err)
	}

	// -------------------- Administrative Account --------------------
	jlog.Println("You must create an administrative account.")
	var username string
	var password string
	var firstName string

	jlog.Printf("Enter a username:")
	fmt.Scanln(&username)
	jlog.Printf("Enter a password:")
	fmt.Scanln(&password)
	jlog.Printf("Enter First Name:")
	fmt.Scanln(&firstName)

	jlog.Println(username, password)
	adminPassHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		removeDatabase()
		jlog.Errorf("FATAL ERROR: Unable to hash administrative password.\n")
		jlog.Errorf("%v\n", err)
	}

	jlog.Println(string(adminPassHash))

	adminSQL := `INSERT INTO users (username, password, firstName, permissions
) VALUES (?, ?, ?, ?)`
	statement, err := newDB.Prepare(adminSQL)
	if err != nil {
		removeDatabase()
		jlog.Errorf("FATAL ERROR: SQL statement for admin account incorrect.\n")
		jlog.Errorf("%v\b", err)
	}

	_, err = statement.Exec(username, adminPassHash, firstName, 2)
	if err != nil {
		removeDatabase()
		jlog.Errorf("FATAL ERROR: Unable to insert administrative information into database.\n")
		jlog.Errorf("%v\n", err)
	}
	// -------------------- END Administrative Account --------------------

	// -------------------- Login Sessions --------------------
	jlog.Println("Creating login sessions table...")

	loginSQL := `CREATE TABLE loginSessions (id INTEGER PRIMARY KEY, cookie TEXT NOT NULL, username TEXT NOT NULL, firstName TEXT NOT NULL, permissions INTEGER NOT NULL, creationTime INTEGER NOT NULL)`
	_, err = newDB.Exec(loginSQL)	
	if err != nil {
		removeDatabase()
		jlog.Errorf("FATAL ERROR: Unable to create login sessions table.\n")
		jlog.Errorf("%v\n", err)
	}
	// -------------------- END Login Sessions --------------------
}

func removeDatabase() {
	os.Remove("./data/jablko.db")	
}

func (instance *JablkoDB) AddUser(username string, password string, firstName string, permissions int) error {
	userSQL := `INSERT INTO users (username, password, firstName, permissions) VALUES(?, ?, ?, ?)`

	statement, err := instance.Db.Prepare(userSQL)
	if err != nil {
		jlog.Errorf("Error in preparing user create SQL statement\n")
		return err
	}

	_, err = statement.Exec(username, password, firstName, permissions)
	if err != nil {
		jlog.Errorf("Error inserting new user in database.\n")
		return err
	}

	return nil
}

func (instance *JablkoDB) AuthenticateUser(username string, password string) (bool, types.UserData) {
	statement, err := instance.Db.Prepare("SELECT * FROM users WHERE username=(?)")
	if err != nil {
		jlog.Errorf("ERROR: Authenticate user SQL is invalid.\n")
	}

	res, err := statement.Query(username)
	if err != nil {
		jlog.Errorf("ERROR: Unable to retrieve user data.\n")
	}
	defer res.Close()

	var authenticated = false

	user := types.UserData{}

	for res.Next() {
		
		err = res.Scan(&user.Id, &user.Username, &user.Password, &user.FirstName, &user.Permissions)

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

		if err == nil {
			authenticated = true
			break
		}
	}

	return authenticated, user
}
