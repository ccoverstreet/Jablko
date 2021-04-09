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
	"sync"

	"github.com/gorilla/mux"
)

type ModManager struct {
	StateMap map[string]ModData
}

type ModData struct {
	RWMutex sync.RWMutex
	Name string `json:"name"`
	Source string `json:"source"`
	Config interface{} `json:"config"`
}

func NewModManager(config string) (*ModManager, error) {
	x := new(ModManager)
	f := make(map[string]ModData)

	err := json.Unmarshal([]byte(config), &f)
	if err != nil {
		return nil, err
	}

	x.StateMap = f
	log.Println(x)

	b, err := json.Marshal(f["test1"])
	if err != nil {
		return nil, err
	}

	log.Println(string(b))

	return x, nil
}

func (mm *ModManager) HandleRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(vars)

	if vars["state"] != "stateless" || vars["state"] == "stateful" {
		log.Printf(`Request "%s" invalid state option "%s"`, r.URL, vars["state"])	
		return 
	}

	stateless := true

	if vars["state"] != "stateless" {
		stateless = false
	}

	mm.passRequest(w, r, stateless)
}

func (mm *ModManager) passRequest(w http.ResponseWriter, r *http.Request, stateless bool) {
	// WLock is called in the modify response portion of
	// the reverse proxy handler. RLock is used on the 
	// initial stateful request since the change of state
	// only occurs after the response comes back.

	log.Println("PASS REQUEST (NOT IMPLEMENTED)", r.URL)	

}
