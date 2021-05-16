// Mod Manager
// Cale Overstreet
// Apr. 24, 2021

// Response for process management for JMODs, passing data to
// JMODs, installing/upgrading JMODS.

package modmanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"

	"github.com/buger/jsonparser"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/ccoverstreet/Jablko/core/jutil"
	"github.com/ccoverstreet/Jablko/core/subprocess"
)

type ModManager struct {
	sync.RWMutex
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

func (mm *ModManager) SaveConfigToFile() {
	mm.Lock()
	defer mm.Unlock()

	log.Info().
		Msg("Saving JMOD data to jmods.json")

	configByte, err := json.Marshal(mm.ProcMap)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Unable to marshal mod manager")
	}

	log.Printf("%s", configByte)
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
	mm.RLock()
	defer mm.RUnlock()

	for key, jmod := range mm.ProcMap {
		if jmod.Port == portNumber && jmod.Key == jmodKey {
			return true, key
		}
	}

	return false, ""
}

func (mm *ModManager) ServiceHandler(w http.ResponseWriter, r *http.Request) {
	// Uses the JMOD-KEY and PORT-NUMBER assigned to each
	// JMOD for authentication. JMODs can save their configs
	// or retrieve information

	log.Printf("ASDASDASD")

	// Check JMOD-KEY header value
	keyValue := r.Header.Get("JMOD-KEY")
	if keyValue == "" {
		log.Error().
			Msg("Empty JMOD-KEY")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Empty JMOD-KEY")
		return
	}

	portValueStr := r.Header.Get("JMOD-PORT")
	portValue, err := strconv.Atoi(portValueStr)
	if portValueStr == "" || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid JMOD-PORT")
		return
	}

	isValid, modName := mm.IsValidService(keyValue, portValue)
	if !isValid {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Specified service doesn't exist")
		return
	}

	vars := mux.Vars(r)

	switch vars["func"] {
	case "saveConfig":
		mm.saveModConfig(w, r, modName)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid function requested")
	}
}

func (mm *ModManager) saveModConfig(w http.ResponseWriter, r *http.Request, modName string) {
	newConfigByte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to read body")
		return
	}

	fmt.Println("ASDJAJSSAJDSJ DMOD ASMDASD")

	mm.ProcMap[modName].Config = newConfigByte

	go mm.SaveConfigToFile()
}
