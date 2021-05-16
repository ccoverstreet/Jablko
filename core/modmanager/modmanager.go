// Mod Manager
// Cale Overstreet
// Apr. 24, 2021

// Response for process management for JMODs, passing data to
// JMODs, installing/upgrading JMODS.

package modmanager

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/buger/jsonparser"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/ccoverstreet/Jablko/core/jutil"
	"github.com/ccoverstreet/Jablko/core/subprocess"
)

type ModManager struct {
	ProcMap map[string]*subprocess.Subprocess
}

var curPort = 44100

func NewModManager(conf []byte) (*ModManager, error) {
	newMM := new(ModManager)
	newMM.ProcMap = make(map[string]*subprocess.Subprocess)

	// Creates subprocesses for all
	parseConfObj := func(key []byte, value []byte, _ jsonparser.ValueType, _ int) error {
		jmodKey, err := jutil.RandomString(32)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Unable to generate random string for jmodKey")

			panic(err)
		}

		newMM.ProcMap[string(key)] = subprocess.CreateSubprocess(string(key), 8080, curPort, jmodKey, "./data", value)
		curPort += 1
		return nil
	}

	jsonparser.ObjectEach(conf, parseConfObj)

	// Try to start all subprocesses
	for key, subProc := range newMM.ProcMap {
		err := subProc.Build()
		if err != nil {
			log.Error().
				Err(err).
				Caller().
				Str("JMOD", key).
				Msg("Unable to build JMOD")

			continue
		}

		err = subProc.Start()
		if err != nil {
			log.Error().
				Err(err).
				Str("subprocess", key).
				Msg("Unable to start subprocess")
		}
	}

	return newMM, nil
}

func (mm *ModManager) PassRequest(w http.ResponseWriter, r *http.Request) {
	source := r.FormValue("JMOD-Source")

	modPort := mm.ProcMap[source].Port
	url, _ := url.Parse("http://localhost:" + strconv.Itoa(modPort))
	proxy := httputil.NewSingleHostReverseProxy(url)

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to read incoming proxy request body")
	}

	r.Host = url.Host
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	proxy.ServeHTTP(w, r)
}

func (mm *ModManager) IsValidService(jmodKey string, portNumber int) (bool, string) {
	return true, ""
}

func (mm *ModManager) ServiceHandler(w http.ResponseWriter, r *http.Request) {
	// Uses the JMOD-KEY and PORT-NUMBER assigned to each
	// JMOD for authentication. JMODs can save their configs
	// or retrieve information

	// Check JMOD-KEY header value
	keyValue := r.Header.Get("JMOD-KEY")
	if keyValue == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Empty JMOD-KEY")
		return
	}

	portValue := r.Header.Get("JMOD-PORT")
	if portValue == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Empty JMOD-KEY")
		return
	}

	vars := mux.Vars(r)

	switch vars["func"] {
	case "saveConfig":

	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid function requested")
	}
}
