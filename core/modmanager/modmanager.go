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
	"time"

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

		newMM.ProcMap[string(key)] = subprocess.CreateSubprocess(string(key), 8080, jmodKey, "./data", value)
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

		go subProc.Start()
	}

	return newMM, nil
}

func (mm *ModManager) SaveConfigToFile() error {
	mm.Lock()
	defer mm.Unlock()

	log.Info().
		Msg("Saving JMOD data to jmods.json")

	configByte, err := json.MarshalIndent(mm.ProcMap, "", "    ")
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Unable to marshal mod manager")

		return err
	}

	err = ioutil.WriteFile("./jmods.json", configByte, 0666)

	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Unable to save jmods.json")

		return err
	}

	return nil
}

func (mm *ModManager) PassRequest(w http.ResponseWriter, r *http.Request) {
	source := r.FormValue("JMOD-Source")

	modPort := mm.ProcMap[source].ModPort
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

func (mm *ModManager) JMODData() ([]byte, error) {
	mm.Lock()
	defer mm.Unlock()

	return json.Marshal(mm.ProcMap)
}

func (mm *ModManager) IsJMODStopped(jmodName string) bool {
	mm.Lock()
	defer mm.Unlock()

	if proc, ok := mm.ProcMap[jmodName]; ok {
		if proc.Cmd.Process != nil && proc.Cmd.ProcessState != nil {
			return true
		}
	}

	return false
}

func (mm *ModManager) StartJMOD(jmodName string) error {
	mm.Lock()
	defer mm.Unlock()

	if subProc, ok := mm.ProcMap[jmodName]; ok {
		err := subProc.Start()

		// Check for a three second period if process
		// is still considered as running. This is for
		// handling restarts
		if err != nil {
			if err.Error() == "Process is already started" {
				for i := 0; i < 3; i++ {
					log.Warn().
						Str("jmodName", jmodName).
						Msg("Retrying mod start")

					time.Sleep(1 * time.Second)
					err = subProc.Start()

					if err == nil {
						break
					}
				}
			}
		}

		return err
	}

	return fmt.Errorf("JMOD not found")
}

func (mm *ModManager) StopJMOD(jmodName string) error {
	mm.Lock()
	defer mm.Unlock()

	if subProc, ok := mm.ProcMap[jmodName]; ok {
		return subProc.Stop()
	}

	return fmt.Errorf("JMOD not found")
}

func (mm *ModManager) SetJMODConfig(jmodName string, newConfig string) error {
	mm.Lock()
	defer mm.Unlock()

	if proc, ok := mm.ProcMap[jmodName]; ok {
		proc.Lock()
		defer proc.Unlock()
		proc.Config = []byte(newConfig)

		return nil
	}

	return fmt.Errorf("JMOD not found")
}

func (mm *ModManager) CleanProcesses() {
	for name, proc := range mm.ProcMap {
		log.Info().
			Str("jmodName", name).
			Msg("Cleaning up JMOD process")

		err := proc.Stop()
		if err != nil {
			log.Info().
				Err(err).
				Str("jmodName", name).
				Msg("Unable to clean up JMOD")
		}
	}
}

// ---------- Routes called by JMODs ----------

// Uses the JMOD-KEY and PORT-NUMBER assigned to each
// JMOD for authentication. JMODs can save their configs
// or retrieve information
func (mm *ModManager) ServiceHandler(w http.ResponseWriter, r *http.Request) {

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

func (mm *ModManager) IsValidService(jmodKey string, portNumber int) (bool, string) {
	mm.RLock()
	defer mm.RUnlock()

	for key, jmod := range mm.ProcMap {
		if jmod.ModPort == portNumber && jmod.Key == jmodKey {
			return true, key
		}
	}

	return false, ""
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
