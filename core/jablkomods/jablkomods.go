// Jablko Mod Manager
// Cale Overstreet
// Mar. 30, 2021

// Responsible for managing mod state and jablkomod
// processes. Handles routing related to jmod and
// pmod routes.

package jablkomods

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"io/ioutil"
	"sync"
	"bytes"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/ccoverstreet/Jablko/core/subprocess"
)

type ModManager struct {
	StateMap map[string]*ModData
	SubprocessMap map[string]*subprocess.Subprocess
}

type ModData struct {
	sync.RWMutex
	Name string `json:"name"`
	Source string `json:"source"`
	Config interface{} `json:"config"`
}

func NewModManager(config string) (*ModManager, error) {
	mm := new(ModManager)
	mm.SubprocessMap = make(map[string]*subprocess.Subprocess)
	f := make(map[string]*ModData)

	err := json.Unmarshal([]byte(config), &f)
	if err != nil {
		return nil, err
	}

	mm.StateMap = f
	log.Println(mm)

	b, err := json.Marshal(f["test1"])
	if err != nil {
		return nil, err
	}

	log.Println(string(b))

	return mm, nil
}

func (mm *ModManager) StartJablkoMod(source string) error {
	newProc, err := subprocess.CreateSubprocess(
		source,
		8080,
		8081,
		"./data",
	)

	if err != nil {
		panic(err)
		return err
	}

	log.Printf(`Building "%s"...`, source)
	err = newProc.Build()
	if err != nil {
		panic(err)
		return err
	}

	err = newProc.Start()
	if err != nil {
		return err
	}

	mm.SubprocessMap[source] = newProc

	return nil
}

func (mm *ModManager) HandleRequest(w http.ResponseWriter, r *http.Request) {
	// Checks given parameters to see if valid
	// values are provided

	log.Println("MOD MANAGER HANDLE REQUEST")
	vars := mux.Vars(r)

	modSource := ""

	// Check if modId is in StateMap
	// Send 404 error if not
	if val, ok := mm.StateMap[vars["modId"]]; !ok {
		http.Error(w, "Mod not found.", http.StatusNotFound)
		return
	} else {
		modSource = val.Source
	}

	if _, ok := mm.SubprocessMap[modSource]; !ok{
		http.Error(w, "Subprocess not found.", http.StatusNotFound)
		return
	}

	if vars["state"] != "stateless" && vars["state"] != "stateful" {
		log.Printf(`Request "%s" invalid state option "%s"`, r.URL, vars["state"])
		return
	}

	stateless := true

	if vars["state"] != "stateless" {
		stateless = false
	}

	if r.Method == "GET" {
		mm.passRequest(w, r, vars["modId"], modSource, stateless, true)
	} else {
		mm.passRequest(w, r, vars["modId"], modSource, stateless, false)
	}
}

func (mm *ModManager) passRequest(w http.ResponseWriter, r *http.Request, modId string, modSource string, stateless bool, isGET bool) {
	// WLock is called in the modify response portion of
	// the reverse proxy handler. RLock is used on the 
	// initial stateful request since the change of state
	// only occurs after the response comes back.

	// Get the mod state to read instance data
	// Time holding this lock isn't too
	// relevant since multiple readers are
	// allowed.
	modState := mm.StateMap[modId]
	modState.RLock()
	defer modState.RUnlock()

	modPort := mm.SubprocessMap[modSource].Port

	url, _ := url.Parse("http://localhost:" + strconv.Itoa(modPort))
	proxy := httputil.NewSingleHostReverseProxy(url)

	sentBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	log.Println(string(sentBody))

	if isGET {
		r.Host = url.Host
		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
	} else {
		// Prep Request for proxy
		r.Host = url.Host
		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Header.Set("Content-Type", "application/json")
		//r.Header.Set("MY-SPECIAL-HEADER", "JABLKO SECRET")

		r.Body = ioutil.NopCloser(bytes.NewBuffer(sentBody))
	}

	// Add err handler for proxy
	proxy.ErrorHandler = mm.proxyErrHandler

	// Add res handler for merging state changes
	if !stateless {
		proxy.ModifyResponse = mm.proxyResHandler
	}

	proxy.ServeHTTP(w, r)
}

func (mm *ModManager) proxyResHandler(res *http.Response) error {
	log.Println("PROXY RES HANDLER", res)

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	res.Body = ioutil.NopCloser(bytes.NewBuffer(resData))

	return nil
}

func (mm *ModManager) proxyErrHandler(w http.ResponseWriter, r *http.Request, err error) {
	// This should handle errors from contacting the proxy
	// If error is "connection refused", the status of the
	// subprocess should be check and restarted.

	if strings.Contains(err.Error(), "connection refused") {
		// Connection refused error

		log.Println("PROCESS MIGHT NEED TO BE RESTARTED")
		log.Println(err.Error())
	} else {
		log.Println(err.Error())
	}
}
