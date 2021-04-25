// Jablko Core App
// Cale Overstreet
// Mar. 30, 2021

// Describes how the functionality of Jablko integrate
// into a single struct that is created in the main 
// function.

package app

import (
	"fmt"
	"net/http"
	//"encoding/json"
	"io/ioutil"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/buger/jsonparser"

	"github.com/ccoverstreet/Jablko/core/jablkomods"
	"github.com/ccoverstreet/Jablko/core/modmanager"
)

type JablkoCoreApp struct {
	Router *mux.Router
	ModManager *jablkomods.ModManager
	ModM *modmanager.ModManager
}

func (app *JablkoCoreApp) Init() error {
	// Runs through procedures to instantiate
	// config data.
	if err := app.initRouter(); err != nil {
		return err
	}

	// Read jablkoconfig.json
	confByte, err := ioutil.ReadFile("./jablkoconfig.json")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to read jablkoconfig.json")

		return err
	}


	sourceConf, _, _, err := jsonparser.Get(confByte, "sources")
	if err != nil {
		panic(err)
	}
	log.Printf("%s", sourceConf)

	newModM, err := modmanager.NewModManager(sourceConf)
	if err != nil {
		panic(err)
	}
	log.Printf("%v", newModM)
	app.ModM = newModM

	// jablkomods WILL BE REMOVED
	/*
	newManager, err := jablkomods.NewModManager(string(confByte))

	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to create ModManager")
		return err
	}
	app.ModManager = newManager
	*/

	return nil
}

func (app *JablkoCoreApp) initRouter() error {
	// Creates the gorilla/mux router passed to 
	// http.ListenAndServe

	router := mux.NewRouter()
	router.HandleFunc("/", app.DashboardHandler).Methods("GET")
	router.HandleFunc("/{client}/{func}", app.PassToJMOD).Methods("GET", "POST")
	//router.HandleFunc("/{client}/{state}/{modId}/{modFunc}", app.PassToModManager).Methods("POST", "GET")

	app.Router = router

	return nil
}

func (app *JablkoCoreApp) PassToJMOD(w http.ResponseWriter, r *http.Request) {
	// Checks for JMOD_Source URL parameter
	// Returns 404
	source := r.FormValue("JMOD_Source")
	log.Printf("Source: '%s' %d", source, len(source))


	// Check if no JMOD-Source header value found
	if len(source) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Empty JMOD_Source parameter")
		log.Warn().
			Str("JMOD-Source", source).
			Msg("Empty JMOD_Source parameter")
		return
	}

	// Check if JMOD-Source is a valid option
	if _, ok := app.ModM.ProcMap[source]; ok {
		app.ModM.PassRequest(w, r)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Haven't implemented this yet")
}

/* THIS WILL BE REMOVED IN THE FUTURE
func (app *JablkoCoreApp) PassToModManager(w http.ResponseWriter, r *http.Request) {
	// This wrapper function is needed for a non-nil
	// pointer to be passed to the ModManager 
	// methods.

	app.ModManager.HandleRequest(w, r)
}
*/

func (app *JablkoCoreApp) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	log.Trace().
		Str("reqIPAddress", r.RemoteAddr).
		Msg("Dashboard Handler requested")

	b, err := ioutil.ReadFile("./html/index.html")
	if err != nil {
		return
	}

	fmt.Fprintf(w, "%s", b)

	return
}
