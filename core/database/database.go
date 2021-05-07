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
	//"database/sql"
	"fmt"
	"io/ioutil"
	"sync"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	
	"github.com/rs/zerolog/log"
)

type user struct {
	PasswordHash string `json:"passwordHash"`
}

type pmod struct {
	Key string `json:"key"`
}

type session struct {
	cookieValue string
	creationTime int
}

// userSessions should be indexed by cookie value
// to prevent users from being logged out by 
// overwrites
type DatabaseHandler struct {
	sync.RWMutex
	filePath string
	Users map[string]user `json:"users"`
	Pmods map[string]pmod `json:"pmods"`
	userSessions map[string]session
}

func CreateDatabaseHandler() *DatabaseHandler {
	dh := new(DatabaseHandler)
	dh.Users = make(map[string]user)
	dh.Pmods = make(map[string]pmod)
	dh.userSessions = make(map[string]session)

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

	err = ioutil.WriteFile(db.filePath, b, 0777)

	return nil
}

// Adds user to in memory database and triggers
// a file dump to store state
func (db *DatabaseHandler) CreateUser(username string, password string) error {
	db.Lock()
	defer db.Unlock()

	if _, ok := db.Users[username]; ok {
		return fmt.Errorf("User already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}

	db.Users[username] = user{string(hash)}

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

//func (db *DatabaseHandler) LoginUser
