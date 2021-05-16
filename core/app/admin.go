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
	case "createUser":
		log.Printf("addUser route called")
		app.addUser(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid admin function requested")
	}
}

func (app *JablkoCoreApp) addUser(w http.ResponseWriter, r *http.Request) {
	type submittedData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var data submittedData

	log.Printf("ASDASDASDASDASDASDASDS")
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
