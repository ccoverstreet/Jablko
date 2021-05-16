package database

import (
	"encoding/json"
	"fmt"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	PasswordHash    string `json:"passwordHash"`
	PermissionLevel int    `json:"permissionLevel"`
}

type UserTable struct {
	sync.RWMutex
	Users map[string]user `json:"users"`
}

func CreateUserTable() *UserTable {
	ut := new(UserTable)
	ut.Users = make(map[string]user)

	return ut
}

// This is critical for making the database JSON Unmarshalable
func (ut *UserTable) UnmarshalJSON(data []byte) error {
	ut.Lock()
	defer ut.Unlock()
	err := json.Unmarshal(data, &ut.Users)
	return err
}

// This is critical for making the database JSON Marshalable
func (ut *UserTable) MarshalJSON() ([]byte, error) {
	ut.RLock()
	defer ut.RUnlock()
	b, err := json.Marshal(&ut.Users)
	return b, err
}

func (ut *UserTable) GetUserList() []string {
	ut.RLock()
	defer ut.RUnlock()

	arr := []string{}
	for username, _ := range ut.Users {
		arr = append(arr, username)
	}

	return arr
}

func (ut *UserTable) CreateUser(username string, password string, permissionLevel int) error {
	ut.Lock()
	defer ut.Unlock()

	if len(username) == 0 {
		return fmt.Errorf("Username cannot be empty")
	}

	fmt.Printf("'%s': %d", password, len(password))
	if len(password) < 12 {
		return fmt.Errorf("Password must be at least 12 characters")
	}

	if _, ok := ut.Users[username]; ok {
		return fmt.Errorf("User already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}

	ut.Users[username] = user{string(hash), permissionLevel}

	return nil
}

func (ut *UserTable) DeleteUser(username string) error {
	ut.Lock()
	defer ut.Unlock()

	if _, ok := ut.Users[username]; !ok {
		return fmt.Errorf("User does not exist in database")
	}

	delete(ut.Users, username)

	return nil
}

func (ut *UserTable) IsValidCredentials(username string, password string) (bool, int) {
	ut.RLock()
	defer ut.RUnlock()

	if val, ok := ut.Users[username]; ok {
		res := bcrypt.CompareHashAndPassword([]byte(val.PasswordHash), []byte(password))
		// res == nil on sucess
		if res == nil {
			return true, val.PermissionLevel
		}
		return false, 0
	}

	return false, 0
}
