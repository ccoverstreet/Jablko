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

func NewModManager() *ModManager {
	newMM := &ModManager{sync.RWMutex{}, make(map[string]*subprocess.Subprocess)}

	return newMM
}

func (mm *ModManager) UnmarshalJSON(data []byte) error {
	// Use a map of json.RawMessage as an intermediate
	// This method

	intermediateMap := make(map[string]*json.RawMessage)

	err := json.Unmarshal(data, &intermediateMap)
	if err != nil {
		return err
	}

	for name, rawMessage := range intermediateMap {
		err := mm.AddJMOD(name, *rawMessage)
		if err != nil {
			return err
		}
	}

	return nil
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
	fmt.Println("DELETE", proc, ok)
	if !ok {
		return fmt.Errorf("JMOD not found.")
	}

	err := proc.Stop()
	if err != nil {
		log.Info().
			Err(err).
			Str("jmodName", jmodPath).
			Msg("Process is already stopped")
	}

	if strings.HasPrefix(jmodPath, "github.com") {
		err := github.DeleteSource(jmodPath)
		if err != nil {
			return fmt.Errorf("Unable to remove JMOD source - %v", err.Error())
		}
	}

	delete(mm.ProcMap, jmodPath)

	return mm.SaveConfigToFile()
}

// Does this need to be locked?
// It is only called when the manager is already locked
func (mm *ModManager) SaveConfigToFile() error {
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
			Msg("Unable to save jmods.json")

		return err
	}

	return nil
}

func (mm *ModManager) PassRequest(w http.ResponseWriter, r *http.Request) error {
	mm.RLock()
	defer mm.RUnlock()

	source := r.FormValue("JMOD-Source")

	proc, ok := mm.ProcMap[source]
	if !ok {
		return fmt.Errorf("JMOD does not exist")
	}

	modPort := proc.ModPort
	url, _ := url.Parse("http://localhost:" + strconv.Itoa(modPort))
	proxy := httputil.NewSingleHostReverseProxy(url)

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("Unable to read incoming request body - %v", err)
	}

	r.Host = url.Host
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	proxy.ServeHTTP(w, r)

	return nil

}

type DashComponent struct {
	Err      error
	JMODName string
	WebComp  []byte
	InstConf []byte
}

func (mm *ModManager) GenerateJMODDashComponents() (string, string) {
	mm.RLock()
	defer mm.RUnlock()

	// Create channel and spawn goroutines to query JMODs
	outChan := make(chan DashComponent, 2)
	for jmodName, proc := range mm.ProcMap {
		baseURL := "http://localhost:" + strconv.Itoa(proc.ModPort)
		go getDashComponent(jmodName, baseURL, outChan)
	}

	builderWC := strings.Builder{}
	builderInstance := strings.Builder{}

	// Read from channel the same number of times as
	// the number of JMODs. This is guaranteed to not deadlock as
	// the number of processes (therefore channel reads) is known
	for i := 0; i < len(mm.ProcMap); i++ {
		comp := <-outChan

		if comp.Err != nil {
			log.Error().
				Err(comp.Err).
				Str("jmodName", comp.JMODName).
				Msg("Unable to get dashboard component")

			continue
		}

		builderWC.WriteString("\njablkoWebCompMap[\"" + comp.JMODName + "\"] = ")
		builderWC.Write(comp.WebComp)

		builderInstance.WriteString("\njablkoInstanceConfMap[\"" + comp.JMODName + "\"] = ")
		builderInstance.Write(comp.InstConf)
	}

	return builderWC.String(), builderInstance.String()
}

func getDashComponent(jmodName string, baseURL string, out chan<- DashComponent) {
	bWC, err := queryJMOD(baseURL + "/webComponent")
	if err != nil {
		out <- DashComponent{err, jmodName, nil, nil}
		return
	}

	bID, err := queryJMOD(baseURL + "/instanceData")
	if err != nil {
		out <- DashComponent{err, jmodName, nil, nil}
		return
	}

	out <- DashComponent{nil, jmodName, bWC, bID}
}

func queryJMOD(url string) ([]byte, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Bad status code: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
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

func (mm *ModManager) StartAllJMODs() error {
	for name, _ := range mm.ProcMap {
		err := mm.StartJMOD(name)
		if err != nil {
			return err
		}
	}

	return nil
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

// THIS SHOULD BE MOVED INTO core/app
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

	err = mm.SaveConfigToFile()

	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to save config file")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}
}
