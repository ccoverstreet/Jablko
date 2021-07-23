// Jablko Database Handler
// Cale Overstreet
// May 6, 2021

/*
All accesses to Jablko's database should go through
this handler. Handles creation and state management
for user logins, registered pmods.

Only write processes trigger updates to the database
file. userSessions are not stored on disk.
*/

package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	//"golang.org/x/crypto/bcrypt"
	"github.com/rs/zerolog/log"
)

type pmod struct {
	Key string `json:"key"`
}

// userSessions should be indexed by cookie value
// to prevent users from being logged out by
// overwrites
type DatabaseHandler struct {
	sync.RWMutex
	Users        *UserTable      `json:"users"`
	Pmods        map[string]pmod `json:"pmods"`
	userSessions *SessionsTable

	filePath string
}

func CreateDatabaseHandler() *DatabaseHandler {
	dh := new(DatabaseHandler)
	//dh.Users = make(map[string]user)
	dh.Users = CreateUserTable()

	dh.Pmods = make(map[string]pmod)
	dh.userSessions = CreateSessionsTable()

	return dh
}

func (db *DatabaseHandler) InitEmptyDatabase() {
	username := ""
	fmt.Printf("\nCreating Admin User...\n")

	for {
		fmt.Printf("Enter username for admin user: ")
		_, err := fmt.Scanln(&username)
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
		}

		// Break if username input was accepted
		if username != "" {
			break
		}
	}

	password := ""
	for {
		fmt.Printf("Enter password for admin user: ")
		fmt.Printf("\033[8m") // Makes entered text invisible

		_, err := fmt.Scanln(&password)
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
		}
		fmt.Printf("\033[0m") // Removes no echo formatting

		// Break if password passes
		if len(password) >= 12 {
			break
		}

		fmt.Printf("Password too short\n")
	}

	err := db.CreateUser(username, password, 1)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to create admin user in empty database")

		panic(err)
	}
}

// Loads existing database from JSON file
func (db *DatabaseHandler) LoadDatabase(file string) error {
	db.filePath = file

	// Check if file exists, create if not
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return fmt.Errorf("File does not exist")
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	err = json.Unmarshal(b, db)
	if err != nil {
		return err
	}

	return nil
}

func (db *DatabaseHandler) SaveDatabase() error {
	db.Lock()
	defer db.Unlock()

	b, err := json.MarshalIndent(db, "", "    ")
	if err != nil {
		return err
	}

	log.Printf("%s", b)

	err = ioutil.WriteFile(db.filePath, b, 0666)

	return nil
}

func (db *DatabaseHandler) GetUserList() []string {
	return db.Users.GetUserList()
}

// Adds user to in memory database and triggers
// a file dump to store state
func (db *DatabaseHandler) CreateUser(username string, password string, permissionLevel int) error {
	err := db.Users.CreateUser(username, password, permissionLevel)

	go db.SaveDatabase()

	return err
}

// Delete user from memory database and triggers
// a file dump to store state
func (db *DatabaseHandler) DeleteUser(username string) error {
	err := db.Users.DeleteUser(username)

	go db.SaveDatabase()

	return err
}

func (db *DatabaseHandler) IsValidCredentials(username string, password string) (bool, int) {
	return db.Users.IsValidCredentials(username, password)
}

// Returns the cookie value and an error value
func (db *DatabaseHandler) CreateSession(username string, permissionLevel int) (string, error) {
	return db.userSessions.CreateSession(username, permissionLevel)
}

// Checks database for session
// Purges expired sessions if value is found
// in database but is expired. This is
// kind of a lazy approach to cleaning the database.
// A secondary trigger is needed to prevent the database
// from exploding in size in ideal operating conditions
// (when all cookies are time valid).
func (db *DatabaseHandler) ValidateSession(cookieValue string) (bool, int) {
	return db.userSessions.ValidateSession(cookieValue)
}

func (db *DatabaseHandler) DeleteSession(cookieValue string) {
	db.userSessions.DeleteSession(cookieValue)
}

// HTTP Handler for login route
func (db *DatabaseHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	// Used for JSON Unmarshal
	type userData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid user data")
		return
	}

	var data userData
	err = json.Unmarshal(b, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid user data")
		return
	}

	isValid, permissionLevel := db.IsValidCredentials(data.Username, data.Password)
	if !isValid {
		fmt.Fprintf(w, "Invalid login")
		return
	}

	cookieVal, err := db.CreateSession(data.Username, permissionLevel)
	if err != nil {
		fmt.Fprintf(w, "Unable to create user session")
	}

	cookie := http.Cookie{
		Name:   "jablko-session",
		Value:  cookieVal,
		MaxAge: int(db.userSessions.SessionLifetime),
	}

	http.SetCookie(w, &cookie)
	fmt.Fprintf(w, "success")
}
