package app2

import (
	"fmt"
	"net/http"

	"github.com/ccoverstreet/Jablko/core/jutil"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

var adminFuncMap = map[string]func(*http.Request, *JablkoApp) ([]byte, error){
	"installJMOD": installJMOD,
	"getJMODData": getJMODData,
	"startJMOD":   startJMOD,
	"stopJMOD":    stopJMOD,
	"buildJMOD":   buildJMOD,
	"deleteJMOD":  deleteJMOD,
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
	type installData struct {
		JMODPath string `json:"jmodPath"`
	}

	reqData := installData{}
	err := jutil.ParseJSONBody(r.Body, &reqData)
	if err != nil {
		return nil, err
	}

	err = app.ModM.AddJMOD(reqData.JMODPath, nil)
	if err != nil {
		return nil, err
	}

	err = app.ModM.StartJMOD(reqData.JMODPath)
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
