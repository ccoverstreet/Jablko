package database

import (
	"log"
	"net/http"
	"time"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/ccoverstreet/Jablko/types"
	"github.com/ccoverstreet/Jablko/src/jablkorandom"
)

func CreateSession(database *sql.DB, username string, userData types.UserData) (http.Cookie, error) {
	cookieValue, err := jablkorandom.GenRandomStr(128)
	if err != nil {
		log.Println("ERROR: Unable to generate random string for cookie");
		return http.Cookie{}, err
	}

	statement, err := database.Prepare("INSERT INTO loginSessions (cookie, username, firstName, permissions, creationTime) VALUES (?, ?, ?, ?, strftime('%s', 'now'))")	
	if err != nil {
		log.Println("ERROR: Unable to prepare loginSessions INSERT SQL statement.")
		return http.Cookie{}, err
	}

	_, err = statement.Exec(cookieValue, username, userData.FirstName, userData.Permissions)
	if err != nil {
		log.Println("ERROR: Unable to insert session info into loginSessions")
		return http.Cookie{}, err
	}

	newCookie := http.Cookie {
		Name: "jablkoLogin",
		Value: cookieValue,
		Expires: time.Now().Add(6 * time.Hour),
	}

	return newCookie, nil
}

func DeleteSession(database *sql.DB, cookieValue string) error {
	statement, err := database.Prepare("DELETE FROM loginSessions WHERE cookie=(?)")
	if err != nil {
		return err
	}

	log.Println("Prepared statement")
	_, err = statement.Exec(cookieValue)
	if err != nil {
		return err
	}

	err = CleanSessions(database)
	if err != nil {
		return err
	}

	return nil
}

func CleanSessions(database *sql.DB) error {
	statement, err := database.Prepare("DELETE FROM loginSessions WHERE creationTime < (?)")
	if err != nil {
		return err
	}

	_, err = statement.Exec(time.Now().Unix())
	if err != nil {
		return err
	}

	return nil
}

func ValidateSession(database *sql.DB, cookieValue string) (bool, types.SessionHolder, error) {
	hold := types.SessionHolder{}
	isValid := false

	statement, err := database.Prepare("SELECT * FROM loginSessions WHERE cookie=(?)")		
	if err != nil {
		return false, hold, err
	}

	res, err := statement.Query(cookieValue)
	if err != nil {
		return false, hold, err
	}

	for res.Next() {
		
		err = res.Scan(&hold.Id, &hold.Cookie, &hold.Username, &hold.FirstName, &hold.Permissions, &hold.CreationTime)
		log.Println(err)
		if err == nil {
			break
		}
	}

	res.Close()

	if int64(hold.CreationTime + 8) < time.Now().Unix() {
		// If cookie is expired
		// Delete all cookies from table that are expired
		err = DeleteSession(database, cookieValue)

		return false, hold, err
	} else {
		isValid = true
	}

	return isValid, hold, nil
}
