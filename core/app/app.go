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
	"encoding/json"
	"io/ioutil"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/ccoverstreet/Jablko/core/jablkomods"
)

type JablkoCoreApp struct {
	Router *mux.Router
	ModManager *jablkomods.ModManager
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
		log.Error().Msg("Unable to read jablkoconfig.json")
		return err
	}

	newManager, err := jablkomods.NewModManager(string(confByte))

	if err != nil {
		return err
	}
	app.ModManager = newManager

	return nil
}

func (app *JablkoCoreApp) initRouter() error {
	// Creates the gorilla/mux router passed to 
	// http.ListenAndServe

	router := mux.NewRouter()
	router.HandleFunc("/", app.DashboardHandler).Methods("GET")
	router.HandleFunc("/{client}/{state}/{modId}/{modFunc}", app.PassToModManager).Methods("POST", "GET")

	app.Router = router

	return nil
}

func (app *JablkoCoreApp) PassToModManager(w http.ResponseWriter, r *http.Request) {
	// This wrapper function is needed for a non-nil
	// pointer to be passed to the ModManager 
	// methods.

	app.ModManager.HandleRequest(w, r)
}

func (app *JablkoCoreApp) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	log.Trace().
		Str("reqIPAddress", r.RemoteAddr).
		Msg("Dashboard Handler requested")

	b, err := json.MarshalIndent(app.ModManager.StateMap, "", "  ")
	if err != nil {
		return
	}
	fmt.Fprintf(w, "%s", b)

	return
}
