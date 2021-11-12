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
	intermediateMap := make(map[string]subprocess.JMODData)

	err := json.Unmarshal(data, &intermediateMap)
	if err != nil {
		return err
	}

	for name, data := range intermediateMap {
		err := mm.AddJMOD(name, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mm *ModManager) MarshalJSON() ([]byte, error) {
	mm.RLock()
	defer mm.RUnlock()

	return json.Marshal(mm.ProcMap)
}

func (mm *ModManager) AddJMOD(jmodPath string, jmodData subprocess.JMODData) error {
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
			commit, err := github.RetrieveSource(jmodPath, jmodData.Commit) // Retrieves source
			if err != nil {
				return err
			}

			// Change commit to main branch commit
			jmodData.Commit = commit
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
		jmodData)

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

	return nil
}

func (mm *ModManager) PassRequest(w http.ResponseWriter, r *http.Request) error {
	mm.RLock()

	source := r.FormValue("JMOD-Source")

	proc, ok := mm.ProcMap[source]
	if !ok {
		mm.RUnlock()
		return fmt.Errorf("JMOD does not exist")
	}

	modPort := proc.ModPort
	url, _ := url.Parse("http://localhost:" + strconv.Itoa(modPort))
	proxy := httputil.NewSingleHostReverseProxy(url)

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		mm.RUnlock()
		return fmt.Errorf("Unable to read incoming request body - %v", err)
	}

	r.Host = url.Host
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	mm.RUnlock()
	proxy.ServeHTTP(w, r)

	return nil
}

func (mm *ModManager) GetJMODHostname(jmodName string) (string, error) {
	mm.RLock()
	defer mm.RUnlock()

	proc, ok := mm.ProcMap[jmodName]
	if !ok {
		return "", fmt.Errorf("JMOD not found")
	}

	return "localhost:" + strconv.Itoa(proc.ModPort), nil
}

// Requests to JMODs should be sent through this function as
// requests may need authentication in the future.
func (mm *ModManager) SendRequest(jmodName string, r *http.Request) (*http.Response, error) {
	client := http.Client{}
	return client.Do(r)
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
	bWC, err := QueryJMOD(baseURL + "/webComponent")
	if err != nil {
		out <- DashComponent{err, jmodName, nil, nil}
		return
	}

	bID, err := QueryJMOD(baseURL + "/instanceData")
	if err != nil {
		out <- DashComponent{err, jmodName, nil, nil}
		return
	}

	out <- DashComponent{nil, jmodName, bWC, bID}
}

func QueryJMOD(url string) ([]byte, error) {
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

func (mm *ModManager) StartAllJMODs() []error {
	var errs []error
	for name, _ := range mm.ProcMap {
		err := mm.StartJMOD(name)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Unable to start JMOD")
			errs = append(errs, err)
		}
	}

	return errs
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
// Why the fuck was this being locked in the first place?
func (mm *ModManager) SetJMODConfig(jmodName string, newConfig []byte) error {
	mm.Lock()
	defer mm.Unlock()

	proc, ok := mm.ProcMap[jmodName]
	if !ok {
		return fmt.Errorf("JMOD not found")
	}

	proc.SetConfig(newConfig)

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

func (mm *ModManager) IsValidService(portNumber int, jmodKey string) (bool, string) {
	mm.RLock()
	defer mm.RUnlock()

	for key, jmod := range mm.ProcMap {
		if jmod.ModPort == portNumber && jmod.Key == jmodKey {
			return true, key
		}
	}

	return false, ""
}
