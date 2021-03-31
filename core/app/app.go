// Jablko Core App
// Cale Overstreet
// Mar. 30, 2021

// Describes how the functionality of Jablko integrate
// into a single struct that is created in the main 
// function.

package app

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"

	"github.com/gorilla/mux"

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

	newManager, err := jablkomods.NewModManager(`{
	"test1": {
		"name": "TEST 1",
		"source": "github.com/ccoverstreet/TEST1",
		"config": {}
	},
	"test2": {
		"name": "TEST 2",
		"source": "github.com/ccoverstreet/TEST2",
		"config": {
			"special": 30,
			"updateInterval": 60
		}
	}
}
`)
	if err != nil {
		return err
	}
	app.ModManager = newManager

	log.Println(app.ModManager.StateMap["test1"].Name)

	return nil
}

func (app *JablkoCoreApp) initRouter() error {
	router := mux.NewRouter()
	router.HandleFunc("/", app.DashboardHandler)
	
	app.Router = router

	return nil
}


func (app *JablkoCoreApp) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	b, err := json.MarshalIndent(app.ModManager.StateMap, "", "  ")
	if err != nil {
		return
	}
	fmt.Fprintf(w, "%s", b)

	return
}
