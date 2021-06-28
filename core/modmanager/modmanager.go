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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/buger/jsonparser"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/ccoverstreet/Jablko/core/github"
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
		err := newMM.AddJMOD(string(key), value)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Unable to create new JMOD process")

			return nil
		}

		return nil
	}

	jsonparser.ObjectEach(conf, parseConfObj)

	// Try to start all subprocesses
	for _, subProc := range newMM.ProcMap {
		go subProc.Start()
	}

	return newMM, nil
}

func (mm *ModManager) AddJMOD(jmodPath string, config []byte) error {
	mm.Lock()
	defer mm.Unlock()

	// Check if mod is already instantiated
	if _, ok := mm.ProcMap[jmodPath]; ok {
		return fmt.Errorf("Module is already registered")
	}

	if strings.HasPrefix(jmodPath, "github.com") {
		// Check if mod is already installed
		_, err := os.Stat(jmodPath)
		if os.IsNotExist(err) {
			log.Printf("github.com route called, need to check for download and @ syntax")
			err = github.RetrieveSource(jmodPath) // Retrieves source
			if err != nil {
				return err
			}
		}
	}

	jmodKey, err := jutil.RandomString(32)
	if err != nil {
		return err
	}

	splitJMODPath := strings.Split(jmodPath, "/")
	shortName := splitJMODPath[len(splitJMODPath)-1]

	dataDir, err := filepath.Abs("./data/" + shortName)
	if err != nil {
		return err
	}

	newProc, err := subprocess.CreateSubprocess(jmodPath,
		8080,
		jmodKey,
		dataDir,
		config)

	if err != nil {
		return err
	}

	mm.ProcMap[jmodPath] = newProc

	return nil
}

func (mm *ModManager) BuildJMOD(jmodPath string) error {
	mm.Lock()
	defer mm.Unlock()

	proc, ok := mm.ProcMap[jmodPath]
	if !ok {
		return fmt.Errorf("JMOD does not exist")
	}

	return proc.Build()
}

func (mm *ModManager) DeleteJMOD(jmodPath string) error {
	mm.Lock()
	defer mm.Unlock()

	proc, ok := mm.ProcMap[jmodPath]
	if !ok {
		return fmt.Errorf("JMOD not found.")
	}

	err := proc.Stop()
	if err != nil {
		return err
	}

	delete(mm.ProcMap, jmodPath)

	if strings.HasPrefix(jmodPath, "github.com") {
		err := github.DeleteSource(jmodPath)
		return fmt.Errorf("Unable to remove JMOD source - %v", err.Error())
	}

	go mm.SaveConfigToFile()

	return nil
}

// Does this need to be locked?
// It is only called when the manager is already locked
func (mm *ModManager) SaveConfigToFile() error {
	mm.Lock()
	defer mm.Unlock()

	log.Info().
		Msg("Saving JMOD data to jmods.json")

	configByte, err := json.MarshalIndent(mm.ProcMap, "", "    ")
	if err != nil {
		log.Error().
			Err(err).
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
	mm.RLock()
	defer mm.RUnlock()

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
	mm.RLock()
	defer mm.RUnlock()

	return json.Marshal(mm.ProcMap)
}

func (mm *ModManager) IsJMODStopped(jmodName string) bool {
	mm.RLock()
	defer mm.RUnlock()

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

	proc, ok := mm.ProcMap[jmodName]
	if !ok {
		return fmt.Errorf("JMOD not found")
	}

	err := proc.Start()
	if err != nil {
		// Retry starting for three seconds
		// only if the error is that the process is already started

		if err.Error() == "Process is already started" {
			for i := 0; i < 3; i++ {
				log.Warn().
					Str("jmodName", jmodName).
					Msg("Retrying mod start")

				time.Sleep(1 * time.Second)
				err = proc.Start()

				if err == nil {
					break
				}
			}
		}
	}

	return err
}

func (mm *ModManager) StopJMOD(jmodName string) error {
	mm.Lock()
	defer mm.Unlock()

	proc, ok := mm.ProcMap[jmodName]
	if !ok {
		return fmt.Errorf("JMOD not found")
	}

	return proc.Stop()
}

// Should the proc struct handle this
// Feels weird locking the process outside of process struct
func (mm *ModManager) SetJMODConfig(jmodName string, newConfig string) error {
	mm.Lock()
	defer mm.Unlock()

	proc, ok := mm.ProcMap[jmodName]
	if !ok {
		return fmt.Errorf("JMOD not found")
	}

	proc.Lock()
	defer proc.Unlock()
	proc.Config = []byte(newConfig)

	return nil
}

func (mm *ModManager) GetJMODLog(jmodName string) ([]byte, error) {
	mm.RLock()
	defer mm.RUnlock()

	proc, ok := mm.ProcMap[jmodName]
	if !ok {
		return nil, fmt.Errorf("JMOD process does not exist")
	}

	return proc.GetCurLogBytes()
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
	log.Printf("BIG FART")

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
		log.Error().
			Err(err).
			Msg("Value not set")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid JMOD-PORT")
		return
	}

	isValid, modName := mm.IsValidService(keyValue, portValue)
	if !isValid {
		log.Error().
			Msg("JMOD service does not exist")

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
	mm.Lock()
	defer mm.Unlock()

	log.Info().
		Str("jmodName", modName).
		Msg("JMOD requested config save")

	newConfigByte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to read request")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to read body")
		return
	}

	mm.ProcMap[modName].Config = newConfigByte

	go mm.SaveConfigToFile()
}
