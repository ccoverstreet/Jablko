package database

import (
	"sync"
	"time"

	"github.com/ccoverstreet/Jablko/core/jutil"
)

// PermissionLevel is 0 for regular user and 1 for admin
type session struct {
	username        string
	permissionLevel int
	creationTime    int64 // time.Now().Unix()
}

type SessionsTable struct {
	sync.RWMutex
	Sessions map[string]session

	SessionLifetime int64
}

func CreateSessionsTable() *SessionsTable {
	st := new(SessionsTable)
	st.Sessions = make(map[string]session)

	st.SessionLifetime = 3600

	return st
}

func (st *SessionsTable) CreateSession(username string, permissionLevel int) (string, error) {
	st.Lock()
	defer st.Unlock()

	cookieValue, err := jutil.RandomString(32)
	if err != nil {
		return "", err
	}

	st.Sessions[cookieValue] = session{username, permissionLevel, time.Now().Unix()}

	return cookieValue, nil
}

// Returns a bool for session validity and returns
// permission level of found session.
func (st *SessionsTable) ValidateSession(cookieValue string) (bool, int) {
	st.RLock()
	defer st.RUnlock()

	if val, ok := st.Sessions[cookieValue]; ok {
		// Check if session is expired
		if (time.Now().Unix() - val.creationTime) > st.SessionLifetime {
			// Purge old sessions from database
			go st.CleanSessions()
			return false, 0
		}

		return true, val.permissionLevel
	}

	return false, 0
}

func (st *SessionsTable) DeleteSession(cookieValue string) {
	st.Lock()
	defer st.Unlock()

	if _, ok := st.Sessions[cookieValue]; ok {
		delete(st.Sessions, cookieValue)
	}
}

// Removes expired sessions from the database
// Called when a session value queried is found
// to be invalid and once an hour
func (st *SessionsTable) CleanSessions() {
	st.Lock()
	defer st.Unlock()

	checkTime := time.Now().Unix()
	for cookieValue, data := range st.Sessions {
		if (checkTime - data.creationTime) > st.SessionLifetime {
			delete(st.Sessions, cookieValue)
		}
	}
}
