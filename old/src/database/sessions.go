package database

import (
	"net/http"
	"time"

	"github.com/ccoverstreet/Jablko/types"
	"github.com/ccoverstreet/Jablko/src/jablkorandom"
	"github.com/ccoverstreet/Jablko/src/jlog"
)

const sessionLength = 216000 // Session duration in seconds

func (instance *JablkoDB) CreateSession(username string, userData types.UserData) (http.Cookie, error) {
	cookieValue, err := jablkorandom.GenRandomStr(128)
	if err != nil {
		jlog.Errorf("ERROR: Unable to generate random string for cookie\n");
		return http.Cookie{}, err
	}

	statement, err := instance.Db.Prepare("INSERT INTO loginSessions (cookie, username, firstName, permissions, creationTime) VALUES (?, ?, ?, ?, strftime('%s', 'now'))")	
	if err != nil {
		jlog.Errorf("ERROR: Unable to prepare loginSessions INSERT SQL statement.\n")
		return http.Cookie{}, err
	}

	_, err = statement.Exec(cookieValue, username, userData.FirstName, userData.Permissions)
	if err != nil {
		jlog.Errorf("ERROR: Unable to insert session info into loginSessions\n")
		return http.Cookie{}, err
	}

	newCookie := http.Cookie {
		Name: "jablkoLogin",
		Value: cookieValue,
		Expires: time.Now().Add(6 * time.Hour),
	}

	return newCookie, nil
}

func (instance *JablkoDB) DeleteSession(cookieValue string) error {
	statement, err := instance.Db.Prepare("DELETE FROM loginSessions WHERE cookie=(?)")
	if err != nil {
		return err
	}

	_, err = statement.Exec(cookieValue)
	if err != nil {
		return err
	}

	err = instance.CleanSessions()
	if err != nil {
		return err
	}

	return nil
}

func (instance *JablkoDB) CleanSessions() error {
	statement, err := instance.Db.Prepare("DELETE FROM loginSessions WHERE creationTime < (?)")
	if err != nil {
		return err
	}

	_, err = statement.Exec(time.Now().Unix() - sessionLength)
	if err != nil {
		return err
	}

	return nil

}

func (instance *JablkoDB) ValidateSession(cookieValue string) (bool, types.SessionHolder, error) {
	hold := types.SessionHolder{}
	isValid := false

	statement, err := instance.Db.Prepare("SELECT * FROM loginSessions WHERE cookie=(?)")		
	if err != nil {
		return false, hold, err
	}

	res, err := statement.Query(cookieValue)
	if err != nil {
		return false, hold, err
	}

	for res.Next() {
		err = res.Scan(&hold.Id, &hold.Cookie, &hold.Username, &hold.FirstName, &hold.Permissions, &hold.CreationTime)
		if err == nil {
			break
		}
	}

	res.Close()

	if int64(hold.CreationTime + sessionLength) < time.Now().Unix() {
		// If cookie is expired
		// Delete all cookies from table that are expired
		err = instance.DeleteSession(cookieValue)

		return false, hold, err
	} else {
		isValid = true
	}

	return isValid, hold, nil
}
