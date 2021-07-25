package app

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// ---------- Routes called by JMODs ----------

// Uses the JMOD-KEY and PORT-NUMBER assigned to each
// JMOD for authentication. JMODs can save their configs
// or retrieve information

func (app *JablkoCoreApp) ServiceHandler(w http.ResponseWriter, r *http.Request) {
	modName := r.Header.Get("JMOD-NAME")

	vars := mux.Vars(r)

	switch vars["func"] {
	case "saveConfig":
		app.saveModConfig(w, r, modName)
	case "sendMessage":
		app.sendMessage(w, r, modName)
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

func (app *JablkoCoreApp) sendMessage(w http.ResponseWriter, r *http.Request, modName string) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to read body")
		return
	}

	bodyReader := bytes.NewReader(reqBody)

	for _, messagingMod := range app.MessagingMods {
		bodyReader.Seek(0, 0)
		hostname, err := app.ModM.GetJMODHostname(messagingMod)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Unable to read body")

			continue
		}

		newReq, err := http.NewRequest("POST", "http://"+hostname+"/service/sendMessage", bodyReader)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Unable to create HTTP request")
		}

		_, err = app.ModM.SendRequest(messagingMod, newReq)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Unable to send request to service")

			continue
		}
	}
}
