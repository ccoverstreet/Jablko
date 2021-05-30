package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ccoverstreet/Jablko/core/jutil"
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

	vars := mux.Vars(r)

	switch vars["func"] {
	case "getJMODData":
		app.getJMODData(w, r)
	case "startJMOD":
		app.startJMOD(w, r)
	case "stopJMOD":
		app.stopJMOD(w, r)
	case "applyJMODConfig":
		app.applyJMODConfig(w, r)
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

func (app *JablkoCoreApp) getJMODData(w http.ResponseWriter, r *http.Request) {
	data, err := app.ModM.JMODData()
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Unable to get JMOD data")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}

	log.Printf("%s", data)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", data)
}

func (app *JablkoCoreApp) startJMOD(w http.ResponseWriter, r *http.Request) {
	type startData struct {
		JMODName string `json:"jmodName"`
	}

	var reqData startData

	err := jutil.ParseJSONBody(r.Body, &reqData)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}

	err = app.ModM.StartJMOD(reqData.JMODName)

	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Str("source", reqData.JMODName).
			Msg("Unable to start JMOD")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}

	log.Info().
		Str("source", reqData.JMODName).
		Msg("Started JMOD")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Started jmod")
}

func (app *JablkoCoreApp) stopJMOD(w http.ResponseWriter, r *http.Request) {
	type stopData struct {
		JMODName string `json:"jmodName"`
	}

	var reqData stopData

	err := jutil.ParseJSONBody(r.Body, &reqData)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}

	err = app.ModM.StopJMOD(reqData.JMODName)

	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Str("source", reqData.JMODName).
			Msg("Unable to stop JMOD")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}

	log.Info().
		Str("source", reqData.JMODName).
		Msg("Stopped JMOD")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Stopped jmod")
}

func (app *JablkoCoreApp) applyJMODConfig(w http.ResponseWriter, r *http.Request) {
	type jmodConfig struct {
		JMODName  string `json:"jmodName"`
		NewConfig string `json:"newConfig"`
	}

	var newConfig jmodConfig

	err := jutil.ParseJSONBody(r.Body, &newConfig)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Unable to parse JSON body for JMOD config")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to parse JSON body: %v", err)
		return
	}

	if !json.Valid([]byte(newConfig.NewConfig)) {
		log.Error().
			Caller().
			Msg("Invalid JSON string")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid JSON")
		return
	}

	err = app.ModM.SetJMODConfig(newConfig.JMODName, newConfig.NewConfig)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Unable to set config")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to set JMOD config: %v", err)
		return
	}

	// Restart the JMOD so that changes apply
	err = app.ModM.StopJMOD(newConfig.JMODName)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Unable to stop jmod")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to stop JMOD: %v", err)
		return
	}

	err = app.ModM.SaveConfigToFile()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to save config: %v", err)
		return
	}

	err = app.ModM.StartJMOD(newConfig.JMODName)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Unable to start JMOD")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to start JMOD: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Applied new config")
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

	err := jutil.ParseJSONBody(r.Body, &data)

	if err != nil {
		log.Error().
			Caller().
			Err(err).
			Msg("Unable to read request body")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
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

func (app *JablkoCoreApp) deleteUser(w http.ResponseWriter, r *http.Request) {
	type delUserBody struct {
		Username string `json:"username"`
	}

	var reqData delUserBody

	err := jutil.ParseJSONBody(r.Body, &reqData)

	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Unable to read body for admin/deleteUser")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}

	err = app.DBHandler.DeleteUser(reqData.Username)

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
