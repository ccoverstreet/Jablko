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
	"fmt"
	"io/ioutil"
	"sync"
	"net/http"
	"encoding/json"
	"time"
	"math/big"
	"crypto/rand"

	"golang.org/x/crypto/bcrypt"
	"github.com/rs/zerolog/log"
)

type user struct {
	PasswordHash string `json:"passwordHash"`
	PermissionLevel int `json:"permissionLevel"`
}

type pmod struct {
	Key string `json:"key"`
}

// PermissionLevel is 0 for regular user and 1 for admin
type session struct {
	username string
	permissionLevel int
	creationTime int64 // time.Now().Unix()
}

// userSessions should be indexed by cookie value
// to prevent users from being logged out by 
// overwrites
type DatabaseHandler struct {
	sync.RWMutex
	Users map[string]user `json:"users"`
	Pmods map[string]pmod `json:"pmods"`
	userSessions map[string]session

	filePath string
	sessionLifetime int64 // How long in seconds a session is valid
}

func CreateDatabaseHandler() *DatabaseHandler {
	dh := new(DatabaseHandler)
	dh.Users = make(map[string]user)
	dh.Pmods = make(map[string]pmod)
	dh.userSessions = make(map[string]session)

	// TEMPORARY
	// Used for testing. Should be specified in config in the future
	dh.sessionLifetime = 3600

	return dh
}

// Loads existing database from JSON file
func (db *DatabaseHandler) LoadDatabase(file string) error {
	db.filePath = file

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, db)
	if err != nil {
		return err
	}

	log.Printf("%v", db)

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

// Adds user to in memory database and triggers
// a file dump to store state
func (db *DatabaseHandler) CreateUser(username string, password string, permissionLevel int) error {
	db.Lock()
	defer db.Unlock()

	if _, ok := db.Users[username]; ok {
		return fmt.Errorf("User already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}

	db.Users[username] = user{string(hash), permissionLevel}

	go db.SaveDatabase()

	return nil
}

// Delete user from memory database and triggers
// a file dump to store state
func (db *DatabaseHandler) DeleteUser(username string) error {
	db.Lock()
	defer db.Unlock()

	if _, ok := db.Users[username]; !ok {
		return fmt.Errorf("User does not exist in database")
	}

	delete(db.Users, username)

	go db.SaveDatabase()
	return nil
}

func (db *DatabaseHandler) IsValidCredentials(username string, password string) (bool, int) {
	db.RLock()
	defer db.RUnlock()

	if val, ok := db.Users[username]; ok {
		res := bcrypt.CompareHashAndPassword([]byte(val.PasswordHash), []byte(password))
		// res == nil on sucess
		if res == nil {
			return true, val.PermissionLevel
		}
		return false, 0
	}

	return false, 0
}

// Code for generating random string used as cookie value
const cookieChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Function from https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb
func RandomString(n int) (string, error) {
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(cookieChars))))
		if err != nil {
			return "", err
		}
		res[i] = cookieChars[num.Int64()]
	}

	return string(res), nil
}

// Returns the cookie value and an error value
func (db *DatabaseHandler) CreateSession(username string, permissionLevel int) (string, error) {
	db.Lock()
	defer db.Unlock()

	cookieValue, err := RandomString(32)
	if err != nil {
		return "", err
	}

	db.userSessions[cookieValue] = session{username, permissionLevel, time.Now().Unix()}

	log.Printf("%v", db.userSessions)

	return cookieValue, nil
}

// Checks database for session
// Purges expired sessions if value is found
// in database but is expired. This is 
// kind of a lazy approach to cleaning the database.
// A secondary trigger is needed to prevent the database
// from exploding in size in ideal operating conditions 
// (when all cookies are time valid).
func (db *DatabaseHandler) ValidateSession(cookieValue string) (bool, int) {
	db.RLock()
	defer db.RUnlock()

	if val, ok := db.userSessions[cookieValue]; ok {
		// Check if session is expired
		if (time.Now().Unix() - val.creationTime) > db.sessionLifetime {
				// Purge old sessions from database
				go db.CleanSessions()
				return false, 0
		}

		return true, val.permissionLevel
	}


	return false, 0
}

func (db *DatabaseHandler) DeleteSession(cookieValue string)  {
	db.Lock()
	defer db.Unlock()

	if _, ok := db.userSessions[cookieValue]; ok {
		delete(db.userSessions, cookieValue)
	}
}

// Removes expired sessions from the database
// Called when a session value queried is found 
// to be invalid and once an hour
func (db *DatabaseHandler) CleanSessions() {
	db.Lock()
	defer db.Unlock()

	checkTime := time.Now().Unix()
	for cookieValue, data := range db.userSessions {
		if (checkTime - data.creationTime) > db.sessionLifetime {
			delete(db.userSessions, cookieValue)
		}
	}
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
		Name: "jablko-session",
		Value: cookieVal,
		MaxAge: int(db.sessionLifetime),
	}

	http.SetCookie(w, &cookie)
	fmt.Fprintf(w, "success")
}
