package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// ---------- Routes called by JMODs ----------

// Uses the JMOD-KEY and PORT-NUMBER assigned to each
// JMOD for authentication. JMODs can save their configs
// or retrieve information

func (app *JablkoCoreApp) ServiceHandler(w http.ResponseWriter, r *http.Request) {
	// Verify that request has correct JMOD-KEY and JMOD-PORT
	keyValue := r.Header.Get("JMOD-KEY")
	if keyValue == "" {
		log.Error().
			Msg("Empty JMOD-KEY")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Empty JMOD-KEY header")
		return
	}

	portValue, err := strconv.Atoi(r.Header.Get("JMOD-PORT"))
	if err != nil {
		log.Error().
			Err(err).
			Msg("Invalid JMOD-PORT header value")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid JMOD-PORT header value")
		return
	}

	isValid, modName := app.ModM.IsValidService(portValue, keyValue)
	if !isValid {
		log.Error().
			Msg("No JMOD exists with specified JMOD-PORT and JMOD-KEY header")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No JMOD exists with specified JMOD-PORT and JMOD-KEY header")
		return
	}

	vars := mux.Vars(r)

	switch vars["func"] {
	case "saveConfig":
		app.saveModConfig(w, r, modName)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid service function requested")
	}
}

func (app *JablkoCoreApp) saveModConfig(w http.ResponseWriter, r *http.Request, modName string) {
	newConfigByte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to read body")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to read body")
		return
	}

	err = app.ModM.SetJMODConfig(modName, string(newConfigByte))
	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to set config")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to set config")
		return
	}

	err = app.SaveConfig()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to save modmanager config")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to save modmanager config")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Saved config")
}
