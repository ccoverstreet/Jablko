package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Dispatches admin functions based on incoming HTTP requests
//
func (app *JablkoCoreApp) AdminFuncHandler(w http.ResponseWriter, r *http.Request) {
	// First check if user has correct privileges
	permissionLevel, err := strconv.Atoi(r.Header.Get("Jablko-User-Permissions"))
	if err != nil {
		log.Error().
			Err(err).
			Msg("Jablko-User-Permission header is invalid")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Jablko-User-Permission header is invalid")
		return
	}

	if permissionLevel < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Insufficient permissions")
		return
	}

	log.Printf("%d", permissionLevel)
	vars := mux.Vars(r)

	switch vars["func"] {
	case "getUserList":
		app.getUserList(w, r)
	case "createUser":
		app.addUser(w, r)
	case "deleteUser":
		app.deleteUser(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid admin function requested")
	}
}

func (app *JablkoCoreApp) getUserList(w http.ResponseWriter, r *http.Request) {
	userList := app.DBHandler.GetUserList()

	body, err := json.Marshal(userList)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Unable to marshal userList to JSON")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to marshal userList")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", body)
}

func (app *JablkoCoreApp) addUser(w http.ResponseWriter, r *http.Request) {
	type submittedData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var data submittedData

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().
			Caller().
			Err(err).
			Msg("Unable to read request body")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request body")
		return
	}

	err = json.Unmarshal(reqBody, &data)
	if err != nil {
		log.Error().
			Caller().
			Err(err).
			Msg("Unable to read request body")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body format")
		return
	}

	log.Printf("%s", reqBody)
	log.Printf("%v", data)

	err = app.DBHandler.CreateUser(data.Username, data.Password, 0)

	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Unable to create user")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to create user: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Created user")
}

func (app *JablkoCoreApp) deleteUser(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Unable to read body for admin/deleteUser")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unable to read body for admin/deleteUser")
		return
	}

	type delUserBody struct {
		Username string `json:"username"`
	}

	var body delUserBody

	err = json.Unmarshal(reqBody, &body)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Unable to parse body for admin/deleteUser")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unable to parse body for admin/deleteUser")
		return
	}

	err = app.DBHandler.DeleteUser(body.Username)

	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Unable to delete user for admin/deleteUser")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to delete user for admin/deleteUser")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted user")
}
