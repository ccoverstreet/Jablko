package app2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ccoverstreet/Jablko/core/jutil"
	"github.com/ccoverstreet/Jablko/core/subprocess"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

var adminFuncMap = map[string]func(*http.Request, *JablkoApp) ([]byte, error){
	"installJMOD":     installJMOD,
	"getJMODData":     getJMODData,
	"startJMOD":       startJMOD,
	"stopJMOD":        stopJMOD,
	"buildJMOD":       buildJMOD,
	"deleteJMOD":      deleteJMOD,
	"applyJMODConfig": applyJMODConfig,
	"getJMODLog":      getJMODLog,
	"getUserList":     getUserList,
	"createUser":      createUser,
	"deleteUser":      deleteUser,
}

func AdminFuncHandler(w http.ResponseWriter, r *http.Request, app *JablkoApp) {
	pathVars := mux.Vars(r)
	log.Printf("%v\n", pathVars)

	adminFuncName, ok := pathVars["func"]
	if !ok {
		log.Error().
			Msg("No admin function specified")
		return
	}

	adminFunc, ok := adminFuncMap[adminFuncName]
	if !ok {
		log.Error().
			Msg("Invalid admin function specified")
		return
	}

	data, err := adminFunc(r, app)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Error processing admin function request")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}

	fmt.Fprintf(w, "%s", data)
}

func installJMOD(r *http.Request, app *JablkoApp) ([]byte, error) {
	var reqData struct {
		JMODPath string `json:"jmodPath"`
	}

	err := jutil.ParseJSONBody(r.Body, &reqData)
	if err != nil {
		return nil, err
	}

	// Correctly identify root of JMOD name
	splitPath := strings.Split(reqData.JMODPath, "@")
	jmodName := splitPath[0]

	// Run AddJMOD based on whether a default branch is provided
	if len(splitPath) == 1 {
		// This branch should retrieve the latest default branch commit
		err = app.ModM.AddJMOD(jmodName, subprocess.JMODData{"", nil})
	} else {
		err = app.ModM.AddJMOD(jmodName, subprocess.JMODData{splitPath[1], nil})
	}

	if err != nil {
		return nil, err
	}

	// Start newly downloaded JMOD
	err = app.ModM.StartJMOD(jmodName)
	if err != nil {
		return nil, err
	}

	return []byte("Installed JMOD"), nil
}

func getJMODData(r *http.Request, app *JablkoApp) ([]byte, error) {
	return app.ModM.JMODData()
}

func startJMOD(r *http.Request, app *JablkoApp) ([]byte, error) {
	type startData struct {
		JMODName string `json:"jmodName"`
	}

	reqData := startData{}

	err := jutil.ParseJSONBody(r.Body, &reqData)
	if err != nil {
		return nil, err
	}

	err = app.ModM.StartJMOD(reqData.JMODName)
	if err != nil {
		return nil, err
	}

	return []byte("Started JMOD"), nil
}

func stopJMOD(r *http.Request, app *JablkoApp) ([]byte, error) {
	type stopData struct {
		JMODName string `json:"jmodName"`
	}

	reqData := stopData{}

	err := jutil.ParseJSONBody(r.Body, &reqData)
	if err != nil {
		return nil, err
	}

	err = app.ModM.StopJMOD(reqData.JMODName)
	if err != nil {
		return nil, err
	}

	return []byte("Stopped JMOD"), nil
}

func buildJMOD(r *http.Request, app *JablkoApp) ([]byte, error) {
	type buildData struct {
		JMODName string `json:"jmodName"`
	}

	reqData := buildData{}

	err := jutil.ParseJSONBody(r.Body, &reqData)
	if err != nil {
		return nil, err
	}

	err = app.ModM.BuildJMOD(reqData.JMODName)
	if err != nil {
		return nil, err
	}

	return []byte("Built JMOD"), nil
}

func deleteJMOD(r *http.Request, app *JablkoApp) ([]byte, error) {
	type reqFormat struct {
		JMODName string `json:"jmodName"`
	}

	reqData := reqFormat{}

	err := jutil.ParseJSONBody(r.Body, &reqData)
	if err != nil {
		return nil, err
	}

	err = app.ModM.DeleteJMOD(reqData.JMODName)
	if err != nil {
		return nil, err
	}

	return []byte("Deleted JMOD"), nil
}

func applyJMODConfig(r *http.Request, app *JablkoApp) ([]byte, error) {
	type reqFormat struct {
		JMODName  string `json:"jmodName"`
		NewConfig string `json:"newConfig"`
	}

	reqData := reqFormat{}

	err := jutil.ParseJSONBody(r.Body, &reqData)
	if err != nil {
		return nil, err
	}

	if !json.Valid([]byte(reqData.NewConfig)) {
		return nil, fmt.Errorf("Invalid JSON provided for config")
	}

	err = app.ModM.SetJMODConfig(reqData.JMODName, reqData.NewConfig)
	if err != nil {
		return nil, err
	}

	err = app.SaveConfig()
	if err != nil {
		return nil, err
	}

	err = app.ModM.StopJMOD(reqData.JMODName)
	if err != nil {
		return nil, err
	}

	for i := 0; i < 3; i++ {
		if app.ModM.IsJMODStopped(reqData.JMODName) {
			break
		}
		time.Sleep(1 * time.Second)
	}

	err = app.ModM.StartJMOD(reqData.JMODName)
	if err != nil {
		return nil, err
	}

	return []byte("Applied new config to JMOD"), nil
}

func getJMODLog(r *http.Request, app *JablkoApp) ([]byte, error) {
	type reqFormat struct {
		JMODName string `json:"jmodName"`
	}

	reqData := reqFormat{}

	err := jutil.ParseJSONBody(r.Body, &reqData)
	if err != nil {
		return nil, err
	}

	jmodLog, err := app.ModM.GetJMODLog(reqData.JMODName)
	if err != nil {
		return nil, err
	}

	return jmodLog, nil
}

func getUserList(r *http.Request, app *JablkoApp) ([]byte, error) {
	userList := app.DB.GetUserList()

	body, err := json.Marshal(userList)
	if err != nil {
		return nil, err
	}

	return body, err
}

func createUser(r *http.Request, app *JablkoApp) ([]byte, error) {
	type reqFormat struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	reqData := reqFormat{}

	err := jutil.ParseJSONBody(r.Body, &reqData)
	if err != nil {
		return nil, err
	}

	err = app.DB.CreateUser(reqData.Username, reqData.Password, 0)

	if err != nil {
		return nil, err
	}

	return []byte("Created user"), err
}

func deleteUser(r *http.Request, app *JablkoApp) ([]byte, error) {
	type reqFormat struct {
		Username string `json:"username"`
	}

	reqData := reqFormat{}

	err := jutil.ParseJSONBody(r.Body, &reqData)
	if err != nil {
		return nil, err
	}

	err = app.DB.DeleteUser(reqData.Username)
	if err != nil {
		return nil, err
	}

	return []byte("Deleted user"), nil
}
