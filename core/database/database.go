// Jablko Database Handler
// Cale Overstreet
// May 6, 2021

/*
All accesses to Jablko's database should go through
this handler. Handles creation and state management 
for user logins, registered pmods.
*/

package database

import (
	/*
	"database/sql"

	"github.com/rs/zerolog/log"
	*/
)

type user struct {
	Username string `json:"username"`
	PasswordHash string `json:"passwordHash"`
}

type pmod struct {
	Key string
}

type session struct {
	cookieValue string
	creationTime int
}

type DatabaseHandler struct {
	Users map[string]user
	Pmods map[string]pmod
	UserSessions map[string]session
}

func CreateDatabaseHandler() *DatabaseHandler {
	dh := new(DatabaseHandler)
	dh.Users = make(map[string]user)
	dh.Pmods = make(map[string]pmod)
	dh.UserSessions = make(map[string]session)

	return dh
}
