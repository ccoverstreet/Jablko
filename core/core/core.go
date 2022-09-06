package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ccoverstreet/Jablko/core/procmanager"
	"github.com/gorilla/mux"
)

type JablkoCore struct {
	PMan     procmanager.ProcManager `json:"mods"`
	PortHTTP int                     `json:"portHTTP"`
	router   *mux.Router
}

func CreateJablkoCore(config []byte) (*JablkoCore, error) {
	core := &JablkoCore{procmanager.CreateProcManager(), 8080, mux.NewRouter()}

	err := json.Unmarshal(config, core)
	if err != nil {
		return nil, err
	}

	core.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ASDASDASDASDASDASDAS")
	})

	return core, nil
}

func createRouter(core *JablkoCore) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ASDASDASDASDASDASDAS")
	})

	return r
}

func (core *JablkoCore) Start() {
	http.ListenAndServe(":"+strconv.Itoa(core.PortHTTP), core.router)
}
