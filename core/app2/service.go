package app2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/gorilla/mux"
)

var serviceFuncMap = map[string]func(*http.Request, *JablkoApp, string) ([]byte, error){
	"saveConfig": saveConfig,
}

func ServiceFuncHandler(w http.ResponseWriter, r *http.Request, app *JablkoApp) {
	pathVars := mux.Vars(r)
	modname := r.Header.Get("JMOD-NAME")
	log.Printf("%v\n", pathVars)

	serviceFuncName, ok := pathVars["func"]
	if !ok {
		log.Error().
			Msg("No service function specified")
		return
	}

	serviceFunc, ok := serviceFuncMap[serviceFuncName]
	if !ok {
		log.Error().
			Msg("Invalid service function specified")
		return
	}

	data, err := serviceFunc(r, app, modname)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Error processing service function request")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}

	fmt.Fprintf(w, "%s", data)
}

func saveConfig(r *http.Request, app *JablkoApp, modname string) ([]byte, error) {
	newConfigByte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if !json.Valid(newConfigByte) {
		return nil, fmt.Errorf("Invalid JSON")
	}

	err = app.ModM.SetJMODConfig(modname, string(newConfigByte))
	if err != nil {
		return nil, err
	}

	err = app.SaveConfig()
	if err != nil {
		return nil, err
	}

	return []byte("Saved new JMOD config"), nil
}
